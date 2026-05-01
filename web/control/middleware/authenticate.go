package middleware

import (
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/auth"
)

func MustAuthenticate(authenticator *auth.Authenticator, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwt, err := r.Cookie("session")
		if err != nil {
			http.Redirect(w, r, "/authenticate", 302)
			return
		}

		err = authenticator.VerifyJwt(jwt.Value)
		if err != nil {
			http.Redirect(w, r, "/authenticate", 302)
			return
		}

		next.ServeHTTP(w, r)
	})
}
