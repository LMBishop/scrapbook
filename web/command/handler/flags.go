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

func GetFlags(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		return html.FlagsPage("", "", siteName, site.SiteConfig.Flags), nil
	})
}

func PostFlags(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
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
			return html.FlagsPage("", fmt.Errorf("Could not parse form: %w", err).Error(), siteName, site.SiteConfig.Flags), nil
		}

		var flags config.SiteFlag
		if r.FormValue("disable") == "on" {
			flags = flags | config.FlagDisable
		}
		if r.FormValue("tls") == "on" {
			flags = flags | config.FlagTLS
		}
		if r.FormValue("index") == "on" {
			flags = flags | config.FlagIndex
		}
		if r.FormValue("password") == "on" {
			flags = flags | config.FlagPassword
		}
		if r.FormValue("readonly") == "on" {
			flags = flags | config.FlagReadOnly
		}

		site.SiteConfig.Flags = flags
		err = config.WriteSiteConfig(path.Join(site.Path, "site.toml"), site.SiteConfig)
		if err != nil {
			return html.FlagsPage("", fmt.Errorf("Failed to persist flags: %w", err).Error(), siteName, flags), nil
		}

		return html.FlagsPage(fmt.Sprintf("Successfully set flags %s", site.ConvertFlagsToString()), "", siteName, flags), nil
	})
}
