package errs

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"gopkg.in/alioygur/gores.v1"
)

type Response struct {
	Code    uint16      `json:"code"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Handler struct {
	Debug string
}

func (eh *Handler) Handle(w http.ResponseWriter, err error) {
	apperr, ok := errors.Cause(err).(*Error)
	if eh.Debug == "on" {
		var scode = 422
		if ok {
			scode = apperr.HTTPCode
		}
		gores.String(w, scode, fmt.Sprintf("%+v", err))
		return
	}

	if ok {
		res := Response{Code: apperr.Code, Message: apperr.Error()}
		gores.JSON(w, apperr.HTTPCode, res)
		return
	}

	res := Response{Code: InternalServerError, Message: err.Error()}

	gores.JSON(w, http.StatusInternalServerError, res)
	return
}
