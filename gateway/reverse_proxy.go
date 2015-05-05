package gateway

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var ERR_TIMEOUT = []byte("The server, while acting as a gateway or proxy, did not receive a timely response from the upstream server.")

type ReverseProxy struct {
	service   *Service
	Proxy     *httputil.ReverseProxy
	Transport *http.Transport
}

func (rp *ReverseProxy) Director(r *http.Request) {
	target, err := url.Parse(rp.service.Forward)
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
		closerBuffer io.ReadCloser
		err          error
		w            *http.Response
	)

	w, err = rp.Transport.RoundTrip(r)
	if e, ok := err.(*net.OpError); ok {
		if e.Timeout() {
			closerBuffer = ioutil.NopCloser(bytes.NewBuffer(ERR_TIMEOUT))
			w = &http.Response{
				Request:       r,
				StatusCode:    http.StatusGatewayTimeout,
				ProtoMajor:    r.ProtoMajor,
				ProtoMinor:    r.ProtoMinor,
				ContentLength: int64(len(ERR_TIMEOUT)),
				Body:          closerBuffer,
			}
		}
	}

	return w, err
}

func NewReverseProxy(s *Service) *ReverseProxy {
	rp := &ReverseProxy{service: s}
	rp.Transport = &http.Transport{
		Dial:                timeoutDialer(10*time.Second, 10*time.Second),
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	rp.Proxy = &httputil.ReverseProxy{
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
