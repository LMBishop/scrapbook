package server

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/html"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
)

func ServeSite(siteIndex *index.SiteIndex) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		site := siteIndex.GetSiteByHost(r.Host)
		if site == nil {
			w.WriteHeader(http.StatusNotFound)
			html.NotFoundSitePage(r.Host).Render(w)
			return
		}

		if site.SiteConfig.Flags&config.FlagDisable != 0 {
			w.WriteHeader(http.StatusForbidden)
			html.ForbiddenDisabledPage(site.SiteConfig.Host).Render(w)
			return
		}

		if site.SiteConfig.Flags&config.FlagPassword != 0 {
			jwt, err := r.Cookie("session")
			if err != nil {
				goto deny
			}

			err = site.Authenticator.VerifyJwt(jwt.Value)
			if err != nil {
				goto deny
			}

			goto permit

		deny:
			if strings.HasPrefix(r.URL.Path, "/authenticate") {
				goto ask
			}
			http.Redirect(w, r, "/authenticate?redirect="+url.QueryEscape(r.URL.Path), 302)
			return

		ask:
			handleAsk(w, r, site)
			return

		permit:
		}

		site.Handler.ServeHTTP(w, r)
	}
}

func handleAsk(w http.ResponseWriter, r *http.Request, site *site.Site) {
	redirect := r.URL.Query().Get("redirect")

	switch r.Method {
	case "GET":
		html.AuthenticateSitePage("", redirect, site.Name).Render(w)
	case "POST":
		err := r.ParseForm()
		if err != nil {
			html.AuthenticateSitePage(err.Error(), redirect, site.Name).Render(w)
			return
		}

		password := r.Form.Get("password")

		if len(site.SiteConfig.Password) == 0 || subtle.ConstantTimeCompare([]byte(password), []byte(site.SiteConfig.Password)) != 1 {
			html.AuthenticateSitePage("The password is incorrect", redirect, site.Name).Render(w)
			return
		}

		jwt, err := site.Authenticator.NewJwt()
		if err != nil {
			html.AuthenticateSitePage(fmt.Errorf("Failed to create jwt: %w", err).Error(), redirect, site.Name).Render(w)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  "session",
			Value: jwt,

			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			HttpOnly: true,
		})
		http.Redirect(w, r, redirect, 302)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
