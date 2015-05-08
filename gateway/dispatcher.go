package gateway

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/backstage/backstage/api"
	. "github.com/backstage/backstage/gateway/filter"
)

const DEFAULT_TIMEOUT = 10
const ERR_TIMEOUT = "The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server."
const ERR_NOT_FOUND = "The requested resource could not be found but may be available again in the future. "

type Dispatcher struct {
	handler   *ServiceHandler
	proxy     *ReverseProxy
	Transport *http.Transport
}

func (rp *Dispatcher) Director(r *http.Request) {
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

func (rp *Dispatcher) RoundTrip(r *http.Request) (*http.Response, error) {
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
		log.Printf("Error while accessing %s: %s.", r.Header.Get("X-Forwarded-Host"), err.Error())
		w = ErrorResponse(r, api.InternalServerError(err.Error()))
	}

	return w, nil
}

func NewDispatcher(h *ServiceHandler) *Dispatcher {
	rp := &Dispatcher{handler: h}
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
	rp.proxy = &ReverseProxy{
		Director:  rp.Director,
		Transport: rp,
		Filters:   h.filters,
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
