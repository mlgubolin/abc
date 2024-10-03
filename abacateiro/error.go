package application

import (
	"context"
	"errors"
	"fmt"
)

const (
	ErrInternal           = "internal"
	ErrInvalid            = "invalid"
	ErrInvalidCredentials = "invalid_credentials"
	ErrNotFound           = "not_found"
	ErrConflict           = "conflict"
	ErrNotImplemented     = "not_implemented"
	ErrUnauthorized       = "unauthorized"
	ErrNotAllowed         = "not_allowed"
	ErrUnavailable        = "unavailable"
	ERRCONFLICT           = "conflict"
	ERRINTERNAL           = "internal"
	ERRINVALID            = "invalid"
	ERRNOTFOUND           = "not_found"
	ERRNOTIMPLEMENTED     = "not_implemented"
	ERRUNAUTHORIZED       = "unauthorized"
	ERRNOTALLOWED         = "not_allowed"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("abacateiro error: code=%s, message=%s", e.Code, e.Message)
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return ErrInternal
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "Internal error."
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

// ReportError notifies an external service of errors. No-op by default.
var ReportError = func(ctx context.Context, err error, args ...any) {}

// ReportPanic notifies an external service of panics. No-op by default.
var ReportPanic = func(err any) {}
