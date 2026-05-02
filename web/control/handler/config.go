package handler

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	"github.com/LMBishop/scrapbook/web/control/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetConfig(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)

		cfgStr, err := site.ReadSiteConfigFile()

		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Could not read existing site configuration: %w", err).Error(), site.Name, ""), nil
		} else {
			return html.ConfigPage("", "", site.Name, cfgStr), nil
		}
	})
}

func PostConfig(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)

		err := r.ParseForm()
		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Could not parse form: %w", err).Error(), site.Name, ""), nil
		}

		configStr := r.FormValue("config")
		cfg, err := config.SiteConfigFromString(configStr)
		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Failed to parse configuration: %w", err).Error(), site.Name, configStr), nil
		}

		site.Config = &cfg
		err = site.WriteSiteConfigFile(configStr)
		if err != nil {
			return html.ConfigPage("", fmt.Errorf("Failed to persist config: %w", err).Error(), site.Name, configStr), nil
		}

		site.Initialise()
		index.UpdateSiteIndexes()

		return html.ConfigPage(fmt.Sprintf("Successfully updated config for %s", site.Name), "", site.Name, configStr), nil
	})
}
