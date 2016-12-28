package errs

import (
	"fmt"

	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

// Error Codes
const (
	NotImplementedError uint16 = iota
	InternalServerError
	TokenExpired
)

var (
	ErrTokenExpired = New(TokenExpired, 401, "Token expired")
)

func Unauthorized(msg string, args ...interface{}) *Error {
	return New(NotImplementedError, http.StatusUnauthorized, msg, args)
}

func BadRequest(msg string, args ...interface{}) *Error {
	return New(NotImplementedError, http.StatusBadRequest, msg, args)
}

func New(code uint16, hcode int, msg string, args ...interface{}) *Error {
	return &Error{Msg: msg, Code: code, HTTPCode: hcode, Args: args}
}

// IsTokenExpiredErr checks given error is jwt expired token error
func IsTokenExpiredErr(err error) bool {
	verr, ok := err.(*jwt.ValidationError)
	if ok && verr.Errors == jwt.ValidationErrorExpired {
		return true
	}
	return false
}

// Error is application error
type Error struct {
	Msg      string
	Code     uint16
	HTTPCode int
	Args     []interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf(e.Msg, e.Args...)
}

var NewWithStack = errors.Errorf
var Wrap = errors.WithStack
var WrapMsg = errors.Wrapf
var Cause = errors.Cause
