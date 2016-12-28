package handlers

import (
	"app/interfaces/errs"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// decodeR decodes request's body to given interface
func decodeR(r *http.Request, to interface{}) error {
	return errs.Wrap(json.NewDecoder(r.Body).Decode(to))
}

type response struct {
	Result interface{} `json:"result"`
}

func qParam(k string, r *http.Request) string {
	values := r.URL.Query()[k]

	if len(values) != 0 {
		return values[0]
	}

	return ""
}

func muxVarInt(k string, r *http.Request) (int, error) {
	return strconv.Atoi(mux.Vars(r)[k])
}
