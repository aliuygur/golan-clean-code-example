package errs

import (
	"net/mail"

	"github.com/alioygur/is"
)

// CheckEmail checks input whatever email or not
func CheckEmail(v string) error {
	_, err := mail.ParseAddress(v)
	if err != nil {
		return WrapMsg(BadRequest("Invalid email address"), err.Error())
	}
	return nil
}

func CheckPassword(v string) error {
	if is.StringLength(v, 4, 32) {
		return nil
	}
	return Wrap(BadRequest("Invalid password"))
}

func CheckName(v string) error {
	if is.StringLength(v, 2, 32) {
		return nil
	}
	return Wrap(BadRequest("Invalid name"))
}

func CheckRequired(v, field string) error {
	if len(v) != 0 {
		return nil
	}
	return Wrap(BadRequest("the %s field is required", field))
}

func CheckStringLen(v string, min int, max int, field string) error {
	if is.StringLength(v, min, max) {
		return nil
	}
	return Wrap(BadRequest("the %s field length must between %d and %d", field, min, max))
}
