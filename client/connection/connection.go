package connection

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apihub/apihub"
	"github.com/apihub/apihub/api"
)

//go:generate counterfeiter . Connection
type Connection interface {
	Ping() error
	AddService(apihub.ServiceSpec) (apihub.ServiceSpec, error)
	Services() ([]apihub.ServiceSpec, error)
}

type connection struct {
	client *http.Client
}

func New(listenNetwork, listenAddr string) *connection {
	return &connection{
		client: &http.Client{
			Transport: &http.Transport{
				Dial: func(network, addr string) (net.Conn, error) {
					return net.DialTimeout(listenNetwork, listenAddr, 2*time.Second)
				},
			},
		},
	}
}

func (c *connection) Ping() error {
	return c.do(api.Ping, nil, &struct{}{})
}

func (c *connection) AddService(spec apihub.ServiceSpec) (apihub.ServiceSpec, error) {
	var service apihub.ServiceSpec
	if err := c.do(api.AddService, spec, &service); err != nil {
		return apihub.ServiceSpec{}, err
	}

	return service, nil
}

func (c *connection) Services() ([]apihub.ServiceSpec, error) {
	var specs []apihub.ServiceSpec
	if err := c.do(api.ListServices, nil, &specs); err != nil {
		return []apihub.ServiceSpec{}, err
	}

	return specs, nil
}

func (c *connection) handleError(body io.ReadCloser) error {
	var err apihub.ErrorResponse
	if err := json.NewDecoder(body).Decode(&err); err != nil {
		return errors.New("request failed")
	}
	return errors.New(err.Description)
}

func (c *connection) do(route api.Route, body interface{}, res interface{}) error {
	req, err := createRequest("http://api", route, body)
	if err != nil {
		return err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return c.handleError(response.Body)
	}

	return json.NewDecoder(response.Body).Decode(res)
}

func createRequest(host string, route api.Route, body interface{}) (*http.Request, error) {
	r := api.Routes[route]
	url := fmt.Sprintf("%s/%s", host, strings.TrimLeft(r.Path, "/"))

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, url, buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}
