package handler

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/control/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetConfig(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		cfgStr, err := site.ReadSiteConfigFile()

		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Could not read existing site configuration: %w", err).Error(), siteName, ""), nil
		} else {
			return html.ConfigPage("", "", siteName, cfgStr), nil
		}
	})
}

func PostConfig(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
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
			return html.ConfigPage("", fmt.Errorf("Could not parse form: %w", err).Error(), siteName, ""), nil
		}

		configStr := r.FormValue("config")
		cfg, err := config.SiteConfigFromString(configStr)
		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Failed to parse configuration: %w", err).Error(), siteName, configStr), nil
		}

		site.Config = &cfg
		err = site.WriteSiteConfigFile(configStr)
		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Failed to persist config: %w", err).Error(), siteName, configStr), nil
		}

		site.Initialise()
		index.UpdateSiteIndexes()

		return html.ConfigPage(fmt.Sprintf("Successfully updated config for %s", siteName), "", siteName, configStr), nil
	})
}
