package errs

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	errMsg := "error test,  string arg: %s int arg: %d"
	errArg1 := "hello"
	errArg2 := 1
	err := New(NotImplementedError, http.StatusBadRequest, errMsg, errArg1, errArg2)

	if err.Code != NotImplementedError {
		t.Errorf("Expected error code %d, got %d", NotImplementedError, err.Code)
	}

	if err.HTTPCode != http.StatusBadRequest {
		t.Errorf("Expected error http status code %d, got %d", http.StatusBadRequest, err.HTTPCode)
	}

	if err.Msg != errMsg {
		t.Errorf("Expected error message %s, got %s", errMsg, err.Msg)
	}

	if err.Error() != fmt.Sprintf(errMsg, errArg1, errArg2) {
		t.Errorf("Expected error formatted message %s, got %s", fmt.Sprintf(errMsg, errArg1, errArg2), err.Error())
	}
}
