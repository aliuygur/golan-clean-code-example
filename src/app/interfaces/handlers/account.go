package handlers

import (
	"app"
	"net/http"

	"github.com/alioygur/gores"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type accountService interface {
	ChangeEmail(int, string) error
	ChangePassword(int, string) error
}

func NewAccount(srv accountService, eh app.ErrorHandler) *Account {
	return &Account{srv, eh}
}

type Account struct {
	srv accountService
	eh  app.ErrorHandler
}

func (a *Account) SetRoutes(r *mux.Router, mid ...alice.Constructor) {
	h := alice.New(mid...)
	r.Handle("/v1/me", h.ThenFunc(a.me)).Methods("GET")
	r.Handle("/v1/me", h.ThenFunc(a.update)).Methods("PATCH")
}

func (a *Account) me(w http.ResponseWriter, r *http.Request) {
	u := app.UserMustFromContext(r.Context())

	gores.JSON(w, 200, u)
}

func (a *Account) update(w http.ResponseWriter, r *http.Request) {
	f := new(updateMeForm)
	if err := decodeR(r, f); err != nil {
		a.eh.Handle(w, err)
		return
	}

	me := app.UserMustFromContext(r.Context())

	err := func() error {
		if f.Email != "" {
			if err := a.srv.ChangeEmail(me.ID, f.Email); err != nil {
				return err
			}
		}
		if f.Password != "" {
			if err := a.srv.ChangePassword(me.ID, f.Password); err != nil {
				return err
			}
		}
		return nil
	}()

	if err != nil {
		a.eh.Handle(w, err)
		return
	}

	gores.NoContent(w)
}

type updateMeForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
