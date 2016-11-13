package cfhttp

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"code.cloudfoundry.org/cfhttp/unix_transport"
)

var config Config

type Config struct {
	Timeout time.Duration
}

func Initialize(timeout time.Duration) {
	atomic.StoreInt64((*int64)(&config.Timeout), int64(timeout))
}

func NewClient() *http.Client {
	return newClient(5*time.Second, 0*time.Second, time.Duration(atomic.LoadInt64((*int64)(&config.Timeout))))
}

func NewUnixClient(socketPath string) *http.Client {
	return &http.Client{
		Transport: unix_transport.NewWithDial(socketPath,
			(&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 0 * time.Second,
			}).Dial),
		Timeout: time.Duration(atomic.LoadInt64((*int64)(&config.Timeout))),
	}
}

func NewCustomTimeoutClient(customTimeout time.Duration) *http.Client {
	return newClient(5*time.Second, 0*time.Second, customTimeout)
}

func NewStreamingClient() *http.Client {
	return newClient(5*time.Second, 30*time.Second, 0*time.Second)
}

func newClient(dialTimeout, keepAliveTimeout, timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   dialTimeout,
				KeepAlive: keepAliveTimeout,
			}).Dial,
		},
		Timeout: timeout,
	}
}

func NewTLSConfig(certFile, keyFile, caCertFile string) (*tls.Config, error) {
	tlsCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load keypair: %s", err.Error())
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{tlsCert},
		InsecureSkipVerify: false,
		ClientAuth:         tls.RequireAndVerifyClientCert,
	}

	if caCertFile != "" {
		certBytes, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			return nil, fmt.Errorf("failed read ca cert file: %s", err.Error())
		}

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, errors.New("Unable to load caCert")
		}
		tlsConfig.RootCAs = caCertPool
		tlsConfig.ClientCAs = caCertPool
	}

	return tlsConfig, nil
}
