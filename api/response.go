package api

import "encoding/json"

type HTTPResponse struct {
	StatusCode       int    `json:"-"`
	ErrorType        string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
	Payload          string `json:"message,omitempty"`
}

func (h *HTTPResponse) Output() string {
	if h.ErrorType != "" && h.ErrorDescription != "" {
		r, err := json.Marshal(h)
		if err != nil {
			return err.Error()
		}
		return string(r)
	}
	return h.Payload

}
