package requests

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

type Args struct {
	AcceptableCode int
	Body           interface{}
	Path           string
	Method         string
	Headers        http.Header
}

func (c *HTTPClient) MakeRequest(args Args) (http.Header, int, []byte, error) {
	header := make(map[string][]string)

	url, err := url.Parse(c.Host)
	if err != nil {
		return header, 0, []byte{}, NewInvalidHostError(err)
	}
	url.Path = args.Path

	body, ok := args.Body.(string)
	if !ok {
		body = ""
	}

	req, err := http.NewRequest(args.Method, url.String(), strings.NewReader(body))
	if err != nil {
		return header, 0, []byte{}, NewRequestError(err)
	}

	req.Header = args.Headers
	resp, err := c.client.Do(req)
	if err != nil {
		return header, 0, []byte{}, NewRequestError(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return resp.Header, resp.StatusCode, []byte{}, NewResponseError(err)
	}

	if resp.StatusCode == args.AcceptableCode {
		return resp.Header, resp.StatusCode, respBody, nil
	}

	return resp.Header, resp.StatusCode, respBody, NewResponseError(ErrBadResponse)
}
