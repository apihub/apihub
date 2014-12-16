package api

type HTTPResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
