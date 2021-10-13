package app

import (
	"fmt"
	"runtime"
	"strings"
)

func NewError(msg string) *Error {
	return &Error{
		Msg:  msg,
		Code: strings.ToUpper(strings.ReplaceAll(msg, " ", "_")),
	}
}

func Errorf(format string, args ...interface{}) error {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		panic("unable to get caller info")
	}
	return fmt.Errorf("%s:%d %w", file, line, fmt.Errorf(format, args...))
}

func WrapError(err error) error {
	_, file, line, ok := runtime.Caller(1)

	if !ok {
		panic("unable to get caller info")
	}

	return fmt.Errorf("%s:%d %w", file, line, err)
}

type Error struct {
	Msg  string
	Code string
	Err  error
}

func (e *Error) Error() string {
	return e.Msg
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Wrap(err error) *Error {
	e.Err = err
	return e
}

func (e *Error) GetCode() string {
	return e.Code
}

type ValidationError Error

func NewValidationError(s string) *ValidationError {
	return &ValidationError{Msg: s, Code: "VALIDATION"}
}
