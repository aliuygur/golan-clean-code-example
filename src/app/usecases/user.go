package usecases

import (
	"app"
	"app/interfaces/errs"
	"net/url"
	"os"
	"time"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	errEmailExists    = errs.BadRequest("email address already exists")
	errEmailNotExists = errs.BadRequest("email address not exists")
	errWrongCred      = errs.Unauthorized("wrong creadentials")
	errInactiveUser   = errs.Unauthorized("inactive user")
)

type IUserRepo interface {
	app.Databaser
}

// NewUser instances User
func NewUser(repo IUserRepo, m app.MailSender) *User {
	return &User{repo: repo, m: m}
}

type User struct {
	repo IUserRepo
	m    app.MailSender
}

// Login user
func (u *User) Login(email, password string) (string, error) {
	var usr app.User

	if err := u.repo.OneBy(&usr, app.DBWhere{"Email": email}); err != nil {
		if u.repo.IsNotFoundErr(err) {
			return "", errs.Wrap(errWrongCred)
		}
		return "", err
	}

	if !usr.IsCredentialsVerified(password) {
		return "", errs.Wrap(errWrongCred)
	}

	if !usr.IsActivated {
		return "", errs.Wrap(errInactiveUser)
	}

	token, err := usr.CreateJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}
	return token, nil
}

// Register user
func (u *User) Register(form *RegisterForm) (string, error) {
	if err := form.IsValid(); err != nil {
		return "", err
	}
	// check for email
	exists, err := u.userExistByEmail(form.Email)
	if err != nil {
		return "", err
	} else if exists {
		return "", errs.Wrap(errEmailExists)
	}

	var usr app.User
	usr.FirstName = form.FirstName
	usr.LastName = form.LastName
	usr.Email = form.Email
	usr.IsActivated = true
	usr.SetPassword(form.Password)

	if err := u.repo.Store(&usr); err != nil {
		return "", err
	}

	token, err := usr.CreateJWT(os.Getenv("SECRET_KEY"))
	if err != nil {
		return "", err
	}
	return token, nil
}

// ChangeEmail changes user's email address
func (u *User) ChangeEmail(id int, email string) error {
	if err := errs.CheckEmail(email); err != nil {
		return err
	}

	// check email already exists
	exists, err := u.userExistByEmail(email)
	if err != nil {
		return err
	} else if exists {
		return errs.Wrap(errEmailExists)
	}

	var usr app.User
	usr.ID = id
	return u.repo.UpdateField(&usr, "Email", email)
}

// ChangePassword changes user's password
func (u *User) ChangePassword(id int, pass string) error {
	if err := errs.CheckPassword(pass); err != nil {
		return err
	}

	var usr app.User
	usr.ID = id
	usr.SetPassword(pass)
	return u.repo.UpdateField(&usr, "Password", pass)
}

// SendPasswordResetLink sends to user the password reset link
func (u *User) SendPasswordResetLink(email, resetLink string) error {
	if resetLink == "" {
		resetLink = os.Getenv("PASSWORD_RESET_URL")
	}

	rurl, err := url.Parse(resetLink)
	if err != nil {
		return errs.Wrap(errs.BadRequest("invalid url"))
	}

	// check for email
	exists, err := u.userExistByEmail(email)
	if err != nil {
		return err
	} else if !exists {
		return errs.Wrap(errEmailNotExists)
	}

	claims := jwt.MapClaims{"email": email, "exp": time.Now().Add(time.Hour * 5).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return errs.WrapMsg(err, "token can't signed")
	}

	uquery := rurl.Query()
	uquery.Set("token", tokenString)
	rurl.RawQuery = uquery.Encode()

	body := fmt.Sprintf("Please click below link to reset your password <br/> %s", rurl.String())

	if err := u.m.Send([]string{email}, "şifremi unuttum", []byte(body)); err != nil {
		return err
	}

	return nil
}

func (u *User) ResetPassword(tokenStr, newPassword string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return errs.WrapMsg(err, "token can't parsed, token: %s,", tokenStr)
	}

	email, ok := token.Claims.(jwt.MapClaims)["email"]
	if !ok {
		return errs.NewWithStack("email can't get from token claims, token: %s", tokenStr)
	}

	var usr app.User
	if err := u.repo.OneBy(&usr, app.DBWhere{"Email": email}); err != nil {
		return err
	}

	usr.SetPassword(newPassword)

	if err := u.repo.UpdateField(&usr, "Password", usr.Password); err != nil {
		return err
	}
	return nil
}

func (u *User) userExistByEmail(email string) (bool, error) {
	return u.repo.ExistsBy(&app.User{}, app.DBWhere{"Email": email})
}

type UpdateUserForm struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RegisterForm struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsActivated bool   `json:"isActivated"`
}

func (f *RegisterForm) IsValid() error {
	if err := errs.CheckName(f.FirstName); err != nil {
		return err
	}
	if err := errs.CheckName(f.LastName); err != nil {
		return err
	}
	if err := errs.CheckEmail(f.Email); err != nil {
		return err
	}
	if err := errs.CheckPassword(f.Password); err != nil {
		return err
	}
	return nil
}
