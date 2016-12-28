package usecases

import (
	"app/interfaces/errs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"app"
)

func (u *User) RegisterFacebook(accessToken string) (string, error) {
	fbuser, err := getUserByAccessToken(accessToken)
	if err != nil {
		return "", err
	}

	var usr app.User
	if err := u.repo.OneBy(&usr, app.DBWhere{"Email": fbuser.Email}); err == nil {
		goto token
	} else if !u.repo.IsNotFoundErr(err) {
		return "", err
	}

	usr.FirstName = fbuser.FirstName
	usr.LastName = fbuser.LastName
	usr.Email = fbuser.Email
	usr.IsActivated = true

	if err := u.repo.Store(&usr); err != nil {
		return "", err
	}

token:
	token, err := usr.CreateJWT(os.Getenv("SECRET_KEY"))

	return token, err
}

func getUserByAccessToken(accessToken string) (*fbUser, error) {
	fields := "email,first_name,last_name"
	baseURL := fmt.Sprintf("https://graph.facebook.com/v2.7/me/?access_token=%s&fields=%s", accessToken, fields)

	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	if resp.StatusCode > 299 || resp.StatusCode < 199 {
		return nil, errs.NewWithStack("user can't get from facebook, http status: %d, body: %s", resp.StatusCode, body)
	}

	u := new(fbUser)
	if err := json.Unmarshal(body, u); err != nil {
		return nil, errs.Wrap(err)
	}

	if u.Email == "" {
		return nil, errs.Wrap(errs.BadRequest("there is no email"))
	}

	return u, nil
}

type fbUser struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
