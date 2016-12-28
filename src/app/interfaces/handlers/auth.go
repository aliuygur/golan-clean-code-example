package handlers

import (
	"app"
	"app/usecases"
	"net/http"

	"github.com/gorilla/mux"

	gores "gopkg.in/alioygur/gores.v1"
)

// AuthService interface
type AuthService interface {
	Login(string, string) (string, error)
	Register(*usecases.RegisterForm) (string, error)
	SendPasswordResetLink(string, string) error
	ResetPassword(string, string) error
	RegisterFacebook(string) (string, error)
}

// NewAuthHandler instances new auth handler struct
func NewAuthHandler(s AuthService, eh app.ErrorHandler) *AuthHandler {
	return &AuthHandler{srv: s, eh: eh}
}

// AuthHandler struct
type AuthHandler struct {
	srv AuthService
	eh  app.ErrorHandler
}

// SetRoutes sets this module's routes
func (ah *AuthHandler) SetRoutes(r *mux.Router) {
	r.HandleFunc("/v1/auth/login", ah.login).Methods("POST")
	r.HandleFunc("/v1/auth/register", ah.register).Methods("POST")
	r.HandleFunc("/v1/auth/register-fb", ah.registerFacebook).Methods("POST")

	r.HandleFunc("/v1/password/forgot", ah.forgotPassword).Methods("POST")
	r.HandleFunc("/v1/password/reset", ah.resetPassword).Methods("POST")
}

func (ah *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	f := new(loginform)
	if err := decodeR(r, f); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	token, err := ah.srv.Login(f.Email, f.Password)
	if err != nil {
		ah.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 200, tokenres(token))
}

func (ah *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	f := new(usecases.RegisterForm)
	if err := decodeR(r, f); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	token, err := ah.srv.Register(f)
	if err != nil {
		ah.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 201, tokenres(token))
}

func (ah *AuthHandler) forgotPassword(w http.ResponseWriter, r *http.Request) {
	f := new(forgotPasswordForm)
	if err := decodeR(r, f); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	if err := ah.srv.SendPasswordResetLink(f.Email, f.Link); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	gores.NoContent(w)
}

func (ah *AuthHandler) resetPassword(w http.ResponseWriter, r *http.Request) {
	f := new(resetPasswordForm)
	if err := decodeR(r, f); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	if err := ah.srv.ResetPassword(f.Token, f.Password); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	gores.NoContent(w)
}

func (ah *AuthHandler) registerFacebook(w http.ResponseWriter, r *http.Request) {
	f := new(registerFacebook)
	if err := decodeR(r, f); err != nil {
		ah.eh.Handle(w, err)
		return
	}

	token, err := ah.srv.RegisterFacebook(f.AccessToken)
	if err != nil {
		ah.eh.Handle(w, err)
		return
	}

	gores.JSON(w, 201, tokenres(token))
}

func tokenres(token string) interface{} {
	return struct {
		Token string `json:"token"`
	}{token}
}

type registerFacebook struct {
	AccessToken string `json:"accessToken"`
}

type forgotPasswordForm struct {
	Link  string `json:"link"`
	Email string `json:"email"`
}

type resetPasswordForm struct {
	Password string `json:"password"`
	Token    string `json:"token"`
}

type loginform struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
