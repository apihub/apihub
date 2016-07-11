package connection

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/apihub/apihub/api"
)

type Connection interface {
	Ping() error
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
	return c.do(api.Ping, &struct{}{})
}

func (c *connection) do(route api.Route, res interface{}) error {
	req, err := createRequest("http://api", route)
	if err != nil {
		return err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(res)
}

func createRequest(host string, route api.Route) (*http.Request, error) {
	r := api.Routes[route]
	url := fmt.Sprintf("%s/%s", host, strings.TrimLeft(r.Path, "/"))
	req, err := http.NewRequest(r.Method, url, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
