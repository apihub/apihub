package errors

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}

type HTTPError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Url        string `json:"url"`
}

func (err *HTTPError) Error() string {
	return err.Message
}
