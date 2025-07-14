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

func GetHost(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		return html.HostPage("", "", siteName, site.SiteConfig.Host), nil
	})
}

func PostHost(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
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
			return html.HostPage("", fmt.Errorf("Could not parse form: %w", err).Error(), siteName, site.SiteConfig.Host), nil
		}

		host := r.FormValue("host")

		site.SiteConfig.Host = host
		index.UpdateSiteIndexes()
		err = config.WriteSiteConfig(path.Join(site.Path, "site.toml"), site.SiteConfig)
		if err != nil {
			return html.HostPage("", fmt.Errorf("Failed to persist flags: %w", err).Error(), siteName, host), nil
		}

		return html.HostPage(fmt.Sprintf("Successfully set host for %s to '%s'", siteName, host), "", siteName, host), nil
	})
}
