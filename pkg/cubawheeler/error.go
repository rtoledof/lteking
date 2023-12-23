package cubawheeler

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrAccessDenied      = errors.New("access denied")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInternal          = errors.New("internal server error")
	ErrOrderAccepted     = errors.New("order already accepted")
	ErrNilUserInContext  = errors.New("nil user in context")
	ErrInsufficientFunds = errors.New("insufficient funds")

	ErrInvalid    = errors.New("invalid argument")           // validation failed
	ErrPermission = errors.New("permission denied")          // permission error action cannot be perform.
	ErrExist      = errors.New("already exists")             // entity does exist
	ErrNotExist   = errors.New("does not exist")             // entity does not exist
	ErrConflict   = errors.New("action cannot be performed") // action cannot be performed
)

type Error struct {
	// Human-readable message.
	Message string

	// Param or field with error.
	Param string

	// Underline error.
	Err error

	// Extra data added to the response
	Metadata
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
