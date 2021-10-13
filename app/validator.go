package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Bind(r *http.Request, v interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Errorf("unable to read request body: %w", err)
	}
	if err := json.Unmarshal(data, v); err != nil {
		return Errorf("unable to unmarshal: %s", data)
	}
	return nil
}

func BindAndValidate(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := Bind(r, v); err != nil {
		BadRequest(w, r, err)
		return false
	}

	type validatable interface {
		Validate(r *http.Request) error
	}

	if req, ok := v.(validatable); ok {
		if err := req.Validate(r); err != nil {
			BadRequest(w, r, err)
			return false
		}
	}
	return true
}
