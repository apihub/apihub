package apihub

const (
	E_BAD_REQUEST string = "bad_request"
)

type ErrorResponse struct {
	Error       string `json:"error,omitempty"`
	Description string `json:"error_description,omitempty"`
}
