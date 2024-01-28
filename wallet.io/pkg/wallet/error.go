package wallet

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrForbidden         = errors.New("forbidden")
	ErrAccessDenied      = errors.New("access denied")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternal          = errors.New("internal server error")
	ErrOrderAccepted     = errors.New("order already accepted")
	ErrNilUserInContext  = errors.New("nil user in context")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidCurrency   = errors.New("invalid currency")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrBadRequest        = errors.New("bad request")

	ErrInvalid    = errors.New("invalid argument")           // validation failed
	ErrPermission = errors.New("permission denied")          // permission error action cannot be perform.
	ErrExist      = errors.New("already exists")             // entity does exist
	ErrNotExist   = errors.New("does not exist")             // entity does not exist
	ErrConflict   = errors.New("action cannot be performed") // action cannot be performed
)

type Error struct {
	// Human-readable message.
	Message string `json:"message"`

	// Param or field with error.
	Param string `json:"param,omitempty"`

	StatusCode int `json:"status_code,omitempty"`

	// Underline error.
	Err error `json:"-"`

	// Extra data added to the response
	Metadata Metadata `json:"metadata,omitempty"`
}

func NewError(err error, statusCode int, message string) *Error {
	return &Error{
		Err:        err,
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	var buf bytes.Buffer
	if e.Param != "" {
		fmt.Fprintf(&buf, "Param: %s: ", e.Param)
	}
	buf.WriteString(e.Message)
	if e.Err != nil {
		buf.WriteString(": ")
		buf.WriteString(e.Err.Error())
	}
	return buf.String()
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return (e.Message == t.Message || t.Message == "") &&
		(e.Param == t.Param || t.Param == "")
}

func NewInvalidParameter(param string, value any) *Error {
	return NewError(ErrInvalid, http.StatusBadRequest, fmt.Sprintf("invalid parameter %s: %v", param, value))
}

func NewMissingParameter(param string) *Error {
	return NewError(ErrInvalid, http.StatusBadRequest, fmt.Sprintf("missing parameter %s", param))
}

func NewNotFound(param string) *Error {
	return NewError(ErrNotFound, http.StatusNotFound, fmt.Sprintf("%s not found", param))
}

func NewInternalError(err error) *Error {
	return NewError(err, http.StatusInternalServerError, "internal error")
}
