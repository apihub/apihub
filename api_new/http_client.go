package api_new

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HTTPClient struct {
	Host   string
	client *http.Client
}

func NewHTTPClient(host string) HTTPClient {
	return HTTPClient{
		Host:   host,
		client: &http.Client{},
	}
}

type RequestArgs struct {
	Body   string
	Path   string
	Method string
}

func (c *HTTPClient) MakeRequest(requestArgs RequestArgs) (http.Header, int, []byte, error) {
	url, err := url.Parse(c.Host)
	if err != nil {
		return nil, 0, nil, err
	}

	url.Path = requestArgs.Path
	req, err := http.NewRequest(requestArgs.Method, url.String(), strings.NewReader(requestArgs.Body))
	if err != nil {
		return nil, 0, nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, nil, err
	}

	return resp.Header, resp.StatusCode, respBody, nil
}
