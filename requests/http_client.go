package requests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/apihub/apihub/errors"
	"github.com/franela/goreq"
)

type HTTPClient struct {
	Host string
}

func NewHTTPClient(host string) HTTPClient {
	return HTTPClient{
		Host: host,
	}
}

type Args struct {
	AcceptableCode int
	Body           interface{}
	Path           string
	Method         string
	Headers        http.Header
	Timeout        time.Duration
	ShowDebug      bool
}

func (c *HTTPClient) MakeRequest(args Args) (http.Header, int, []byte, error) {
	header := make(map[string][]string)

	url, err := url.Parse(c.Host)
	if err != nil {
		return header, 0, []byte{}, errors.NewInvalidHostError(err)
	}
	url.Path = args.Path

	req := goreq.Request{
		Uri:       url.String(),
		Method:    args.Method,
		Body:      args.Body,
		Timeout:   args.Timeout,
		ShowDebug: args.ShowDebug,
	}

	for name, value := range args.Headers {
		for _, v := range value {
			req.AddHeader(name, v)
		}
	}

	resp, err := req.Do()
	if err != nil {
		return header, 0, []byte{}, errors.NewRequestError(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.Header, resp.StatusCode, []byte{}, errors.NewResponseError(err)
	}

	if resp.StatusCode == args.AcceptableCode {
		return resp.Header, resp.StatusCode, respBody, nil
	}

	return resp.Header, resp.StatusCode, respBody, errors.NewResponseError(errors.ErrBadResponse)
}
