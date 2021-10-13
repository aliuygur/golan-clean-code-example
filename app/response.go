package app

import (
	"encoding/json"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type HTTPError struct {
	Status int    `json:"status,omitempty"`
	Code   string `json:"code,omitempty"`
	Error  string `json:"error,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, err error, status int) {
	var res = &HTTPError{
		Status: status,
		Code:   "internal",
		Error:  err.Error(),
	}

	// check for error code
	if err, ok := err.(interface {
		GetCode() string
	}); ok {
		res.Code = err.GetCode()
	}

	// check for validation error
	if _, ok := err.(validation.Errors); ok && res.Code == "" {
		res.Code = "validation"
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, err, http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, err, http.StatusNotFound)
}

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, err, http.StatusInternalServerError)
	ReportError(r, err)
}

func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		panic(err)
	}
}
