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
	RemoveService(string) error
	FindService(string) (apihub.ServiceSpec, error)
	UpdateService(string, apihub.ServiceSpec) (apihub.ServiceSpec, error)
}

type Params map[string]string

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
	return c.do(api.Ping, nil, nil, &struct{}{})
}

func (c *connection) AddService(spec apihub.ServiceSpec) (apihub.ServiceSpec, error) {
	var service apihub.ServiceSpec
	if err := c.do(api.AddService, nil, spec, &service); err != nil {
		return apihub.ServiceSpec{}, err
	}

	return service, nil
}

func (c *connection) Services() ([]apihub.ServiceSpec, error) {
	specs := struct {
		Items []apihub.ServiceSpec `json:"items"`
		Count int                  `json:"item_count"`
	}{}

	if err := c.do(api.ListServices, nil, nil, &specs); err != nil {
		return []apihub.ServiceSpec{}, err
	}

	return specs.Items, nil
}

func (c *connection) RemoveService(host string) error {
	params := map[string]string{"host": host}
	return c.do(api.RemoveService, params, nil, &struct{}{})
}

func (c *connection) FindService(host string) (apihub.ServiceSpec, error) {
	params := map[string]string{"host": host}

	var spec apihub.ServiceSpec
	if err := c.do(api.FindService, params, nil, &spec); err != nil {
		return apihub.ServiceSpec{}, err
	}

	return spec, nil
}

func (c *connection) UpdateService(host string, spec apihub.ServiceSpec) (apihub.ServiceSpec, error) {
	params := map[string]string{"host": host}

	var updatedSpec apihub.ServiceSpec
	if err := c.do(api.UpdateService, params, spec, &updatedSpec); err != nil {
		return apihub.ServiceSpec{}, err
	}

	return updatedSpec, nil
}

func (c *connection) hostError(body io.ReadCloser) error {
	var err apihub.ErrorResponse
	if err := json.NewDecoder(body).Decode(&err); err != nil {
		return errors.New("request failed")
	}
	return errors.New(err.Description)
}

func (c *connection) do(route api.Route, params Params, body interface{}, res interface{}) error {
	req, err := createRequest("http://api", route, params, body)
	if err != nil {
		return err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if response.StatusCode < http.StatusOK || response.StatusCode > 299 {
		return c.hostError(response.Body)
	}

	if response.StatusCode != http.StatusNoContent && response.Body != nil {
		defer response.Body.Close()
		return json.NewDecoder(response.Body).Decode(res)
	}

	return nil
}

func createRequest(host string, route api.Route, params Params, body interface{}) (*http.Request, error) {
	r := api.Routes[route]
	url := fmt.Sprintf("%s/%s", host, strings.TrimLeft(r.Path, "/"))

	if params != nil {
		path := strings.TrimLeft(r.Path, "/")
		parts := strings.Split(path, "/")
		for i, p := range parts {
			if p != "" && p[0] == '{' && p[len(p)-1] == '}' {
				param := p[1 : len(p)-1]
				val, ok := params[param]
				if !ok {
					return nil, fmt.Errorf("missing parameter: %s - %s", p, param)
				}
				parts[i] = val
			}
		}
		url = fmt.Sprintf("%s/%s", host, strings.Join(parts, "/"))
		url = strings.TrimLeft(url, "/")
	}

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
