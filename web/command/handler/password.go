package handler

import (
	"fmt"
	"net/http"
	"path"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/command/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetPassword(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		return html.PasswordPage("", "", siteName, site.SiteConfig.Password), nil
	})
}

func PostPassword(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		err := r.ParseForm()
		if err != nil {
			return html.PasswordPage("", fmt.Errorf("Could not parse form: %w", err).Error(), siteName, site.SiteConfig.Password), nil
		}

		password := r.FormValue("password")

		site.SiteConfig.Password = password
		err = config.WriteSiteConfig(path.Join(site.Path, "site.toml"), site.SiteConfig)
		if err != nil {
			return html.PasswordPage("", fmt.Errorf("Failed to persist site: %w", err).Error(), siteName, password), nil
		}

		return html.PasswordPage(fmt.Sprintf("Successfully set password for %s to '%s'", siteName, password), "", siteName, password), nil
	})
}
