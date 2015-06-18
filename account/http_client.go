package account

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/backstage/maestro/errors"
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
	AcceptableCode int
	Body           interface{}
	Path           string
	Method         string
	Headers        http.Header
}

func (c *HTTPClient) MakeRequest(requestArgs RequestArgs) (http.Header, int, []byte, error) {
	header := make(map[string][]string)

	url, err := url.Parse(c.Host)
	if err != nil {
		return header, 0, []byte{}, errors.NewInvalidHostError(err)
	}

	url.Path = requestArgs.Path

	body, ok := requestArgs.Body.(string)
	if !ok {
		body = ""
	}
	req, err := http.NewRequest(requestArgs.Method, url.String(), strings.NewReader(body))
	if err != nil {
		return header, 0, []byte{}, errors.NewRequestError(err)
	}

	req.Header = requestArgs.Headers

	resp, err := c.client.Do(req)
	if err != nil {
		return header, 0, []byte{}, errors.NewRequestError(err)
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.Header, resp.StatusCode, []byte{}, errors.NewResponseError(err)
	}

	if resp.StatusCode == requestArgs.AcceptableCode {
		return resp.Header, resp.StatusCode, respBody, nil
	}

	var errorResponse errors.ErrorResponse
	err = json.Unmarshal(respBody, &errorResponse)
	e := errors.ErrBadResponse
	if err == nil {
		if errorResponse.Description != "" {
			e = errors.NewErrorResponse(errorResponse.Type, errorResponse.Description)
		}
	}
	return resp.Header, resp.StatusCode, respBody, errors.NewResponseError(e)
}
