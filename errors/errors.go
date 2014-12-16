package errors

import "errors"

var (
  ErrUserNotInTeam = errors.New("You do not belong to this team!")
  ErrOnlyOwnerHasPermission = errors.New("Only the owner has permission to perform this operation.")
)

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
