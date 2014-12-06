package api

type HTTPResponse struct {
	StatusCode int    `json:"status_code"`
	Payload    string `json:"message"`
}
