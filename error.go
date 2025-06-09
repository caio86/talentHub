package talenthub

import (
	"errors"
	"fmt"
)

const (
	EINVALID        = "invalid"
	EINTERNAL       = "internal"
	ENOTFOUND       = "not_found"
	ENOTIMPLEMENTED = "not_implemented"
)

type Error struct {
	Code string

	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Error: code=%s message=%s", e.Code, e.Message)
}

func Errorf(code, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "internal error"
}
