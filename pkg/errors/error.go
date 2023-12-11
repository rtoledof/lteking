package errors

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAccessDenied  = errors.New("access denied")
	ErrInvalidInput  = errors.New("invalid input")
	ErrInternal      = errors.New("internal server error")
	ErrOrderAccepted = errors.New("order already accepted")
)
