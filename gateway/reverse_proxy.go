package gateway

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/backstage/backstage/api"
)

const DEFAULT_TIMEOUT = 10
const ERR_TIMEOUT = "The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server."
const ERR_UNEXPECTED_ERROR = "Something went wrong."
const ERR_NOT_FOUND = "The requested resource could not be found but may be available again in the future. "

type ReverseProxy struct {
	handler   *ServiceHandler
	proxy     *httputil.ReverseProxy
	Transport *http.Transport
}

func (rp *ReverseProxy) Director(r *http.Request) {
	target, err := url.Parse(rp.handler.service.Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	targetQuery := target.RawQuery
	r.URL.Scheme = target.Scheme
	r.URL.Host = target.Host
	r.Host = r.URL.Host
	r.URL.Path = joinSlash(target.Path, r.URL.Path)
	if targetQuery == "" || r.URL.RawQuery == "" {
		r.URL.RawQuery = targetQuery + r.URL.RawQuery
	} else {
		r.URL.RawQuery = targetQuery + "&" + r.URL.RawQuery
	}
}

func (rp *ReverseProxy) RoundTrip(r *http.Request) (*http.Response, error) {
	var (
		err error
		w   *http.Response
	)

	w, err = rp.Transport.RoundTrip(r)
	if e, ok := err.(*net.OpError); ok {
		if e.Timeout() {
			w = ErrorResponse(r, api.GatewayTimeout(ERR_TIMEOUT))
		}
	}
	if w == nil && err != nil {
		fmt.Printf("err %+v\n", err)
		w = ErrorResponse(r, api.InternalServerError(ERR_UNEXPECTED_ERROR))
	}

	return w, nil
}

func NewReverseProxy(h *ServiceHandler) *ReverseProxy {
	rp := &ReverseProxy{handler: h}
	t := h.service.Timeout
	if t <= 0 {
		t = DEFAULT_TIMEOUT
	}
	timeout := time.Duration(t)

	rp.Transport = &http.Transport{
		Dial:                timeoutDialer(timeout*time.Second, timeout*time.Second),
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: timeout * time.Second,
	}
	rp.proxy = &httputil.ReverseProxy{
		Director:  rp.Director,
		Transport: rp,
	}
	return rp
}

func timeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func joinSlash(target, path string) string {
	target = strings.TrimSuffix(target, "/")
	path = strings.TrimPrefix(path, "/")
	s := []string{target, path}
	return strings.Join(s, "/")
}

func ErrorResponse(r *http.Request, httpResponse *api.HTTPResponse) *http.Response {
	out := httpResponse.Output()
	var closerBuffer io.ReadCloser = ioutil.NopCloser(bytes.NewBufferString(out))
	w := &http.Response{
		Request:       r,
		StatusCode:    httpResponse.StatusCode,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		ContentLength: int64(len(out)),
		Body:          closerBuffer,
	}
	w.Header = make(map[string][]string)
	w.Header.Add("Content-Type", "application/json")
	return w
}
