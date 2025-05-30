package errs

import (
	"fmt"
)

type CustomError struct {
	Code    ErrorCodeEnum
	Message string
}

func (c *CustomError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", c.Code, c.Message)
}

func New(code ErrorCodeEnum, msg string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: msg,
	}
}

type ErrorCodeEnum int

const (
	CodeUnknown ErrorCodeEnum = iota
	CodeUserNotFound
	CodeInvalidRequest
	CodeParamRequired
)

var (
	ErrUnknown        = New(CodeUnknown, "unknown error")
	ErrInvalidRequest = New(CodeInvalidRequest, "invalid request")
	ErrUserNotFound   = New(CodeUserNotFound, "user not found")
	ErrParamRequired  = New(CodeParamRequired, "user id param required")
)
