package handler

import (
	"crypto/subtle"
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/auth"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/web/command/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetAuthenticate() func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.AuthenticatePage(""), nil
	})
}

func PostAuthenticate(mainConfig *config.MainConfig, authenticator *auth.Authenticator) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			html.AuthenticatePage(err.Error()).Render(w)
			return
		}

		token := r.Form.Get("token")

		if len(mainConfig.Command.Secret) == 0 || subtle.ConstantTimeCompare([]byte(token), []byte(mainConfig.Command.Secret)) != 1 {
			html.AuthenticatePage("The secret key is incorrect").Render(w)
			return
		}

		jwt, err := authenticator.NewJwt()
		if err != nil {
			html.AuthenticatePage(fmt.Errorf("Failed to create jwt: %w", err).Error()).Render(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: jwt,

			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", 302)
	}
}
