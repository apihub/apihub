package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"code.cloudfoundry.org/lager"
)

type transport struct {
	*http.Transport
	logger lager.Logger
}

func (r *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	log := r.logger.Session("round-trip")

	via, err := headerVia(req.Header.Get("Via"), req.ProtoMajor, req.ProtoMinor)
	if err != nil {
		log.Error("failed-read-request-via-hader", err)
		return nil, err
	}
	if via != "" {
		req.Header.Set("Via", via)
	}

	resp, err := r.Transport.RoundTrip(req)
	if err == nil {
		via, err = headerVia(resp.Header.Get("Via"), req.ProtoMajor, req.ProtoMinor)
		if err != nil {
			log.Error("failed-read-response-via-hader", err)
			return nil, err
		}
		if via != "" {
			resp.Header.Set("Via", via)
		}
		return resp, nil
	}

	log.Error("failed-round-trip-request", err)
	respErr := response{
		StatusCode: http.StatusBadGateway,
		Body: responseError{
			ErrType:     "bad_gateway",
			Description: err.Error(),
		},
	}
	if e, ok := err.(*net.OpError); ok {
		if e.Timeout() {
			respErr = response{
				StatusCode: http.StatusGatewayTimeout,
				Body: responseError{
					ErrType:     "gateway_timeout",
					Description: err.Error(),
				},
			}
		}
	}

	return r.Response(req, respErr), nil
}

type response struct {
	StatusCode int
	Body       interface{}
}

type responseError struct {
	ErrType     string `json:"error"`
	Description string `json:"error_description"`
}

func (r *transport) Response(req *http.Request, resp response) *http.Response {
	data, _ := json.Marshal(resp.Body)
	var closerBuffer io.ReadCloser = ioutil.NopCloser(bytes.NewBuffer(data))
	response := &http.Response{
		Request:       req,
		StatusCode:    resp.StatusCode,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		ContentLength: int64(len(data)),
		Body:          closerBuffer,
	}

	response.Header = make(map[string][]string)
	response.Header.Add("Content-Type", "application/json")
	return response
}

func roundTripper(logger lager.Logger, timeout time.Duration) *transport {
	return &transport{
		logger: logger,
		Transport: &http.Transport{
			DialContext:         timeoutDialer(timeout, timeout),
			Proxy:               http.ProxyFromEnvironment,
			TLSHandshakeTimeout: timeout * time.Second,
		},
	}
}

func director(logger lager.Logger, spec ReverseProxySpec) func(req *http.Request) {
	log := logger.Session("create-director")
	log.Debug("start")
	defer log.Debug("end")

	return func(req *http.Request) {
		backend, err := url.Parse(spec.Backends[0])
		if err != nil {
			log.Error("failed-to-parse-backend", err)
			return
		}

		targetQuery := backend.RawQuery
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
		req.Host = req.URL.Host
		backendPath := strings.TrimSuffix(backend.Path, "/")
		reqPath := strings.TrimPrefix(req.URL.Path, "/")
		req.URL.Path = path.Join(backendPath, reqPath)

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
	}
}

func headerVia(original string, protoMajor int, protoMinor int) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	via := strings.Join([]string{original, fmt.Sprintf("%d.%d %s", protoMajor, protoMinor, hostname)}, ", ")
	return strings.Trim(via, ", "), nil
}

func timeoutDialer(timeout time.Duration, rwTimeout time.Duration) func(ctx context.Context, network string, addr string) (net.Conn, error) {
	return func(ctx context.Context, network string, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, addr, timeout)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}
