package interfaces

import (
	"net/http"
	"os"
	"strings"

	"app"

	"app/interfaces/errs"

	"github.com/dgrijalva/jwt-go"
)

type errHandler interface {
	Handle(http.ResponseWriter, error)
}

// getToken gets Authorization key from headers
func getToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errs.Wrap(errs.Unauthorized("invalid authorization token format"))
	}

	return authHeaderParts[1], nil
}

func NewSetUserMid(db app.DBFinder, eh errHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr, err := getToken(r)
			if err != nil {
				eh.Handle(w, err)
				return
			}

			if tokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}

			u, err := func() (*app.User, error) {
				token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("SECRET_KEY")), nil
				})
				if err != nil {
					if errs.IsTokenExpiredErr(err) {
						return nil, errs.Wrap(errs.ErrTokenExpired)
					}
					return nil, errs.WrapMsg(err, "token string can't parsed.")
				}

				userID, ok := token.Claims.(jwt.MapClaims)["userID"].(string)
				if !ok {
					return nil, errs.NewWithStack("userID can't get from token claims, token: %s", tokenStr)
				}

				var u app.User
				if err := db.One(&u, userID); err != nil {
					return nil, err
				}
				return &u, nil
			}()

			if err != nil {
				eh.Handle(w, err)
				return
			}

			ctx := u.NewContext(r.Context())

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func NewAuthRequiredMid(eh errHandler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			usr, ok := app.UserFromContext(r.Context())
			if !ok {
				err := errs.Unauthorized("Auth required")
				eh.Handle(w, errs.Wrap(err))
				return
			}

			if !usr.IsActivated {
				err := errs.Unauthorized("Inactive user")
				eh.Handle(w, errs.Wrap(err))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
