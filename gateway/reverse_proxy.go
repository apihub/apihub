package gateway

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"time"

	"code.cloudfoundry.org/lager"
)

const (
	DEFAULT_TIMEOUT = time.Duration(10)
)

var (
	emptyBackendList = errors.New("Backends cannot be empty.")
)

//go:generate counterfeiter . ReverseProxy
//go:generate counterfeiter . ReverseProxyCreator

type ReverseProxyCreator interface {
	Create(lager.Logger, ReverseProxySpec) (ReverseProxy, error)
}

type ReverseProxy interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ReverseProxySpec struct {
	Handle      string
	Backends    []string
	DialTimeout time.Duration
}

type reverseProxyCreator struct{}

func NewReverseProxyCreator() *reverseProxyCreator {
	return &reverseProxyCreator{}
}

func (rpc *reverseProxyCreator) Create(logger lager.Logger, spec ReverseProxySpec) (ReverseProxy, error) {
	log := logger.Session("reverse-proxy-creator-create")
	log.Info("start", lager.Data{"spec": spec})
	defer log.Info("end")

	if len(spec.Backends) == 0 {
		return nil, emptyBackendList
	}

	timeout := DEFAULT_TIMEOUT
	if spec.DialTimeout > 0 {
		timeout = spec.DialTimeout
	}

	return &reverseProxy{
		spec: spec,
		rp: &httputil.ReverseProxy{
			Director:  director(logger, spec),
			Transport: roundTripper(logger, timeout),
		},
	}, nil
}

type reverseProxy struct {
	spec ReverseProxySpec
	rp   *httputil.ReverseProxy
}

func (n *reverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	n.rp.ServeHTTP(rw, req)
}
