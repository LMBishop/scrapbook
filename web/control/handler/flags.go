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

func GetFlags(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)
		return html.FlagsPage("", "", site.Name, site.Flags), nil
	})
}

func PostFlags(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)

		err := r.ParseForm()
		if err != nil {
			return html.FlagsPage("", fmt.Errorf("Could not parse form: %w", err).Error(), site.Name, site.Flags), nil
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

		site.Flags = flags
		err = site.WriteSiteFlags(flags)
		if err != nil {
			return html.FlagsPage("", fmt.Errorf("Failed to persist flags: %w", err).Error(), site.Name, flags), nil
		}

		site.Initialise()

		return html.FlagsPage(fmt.Sprintf("Successfully set flags %s", site.ConvertFlagsToString()), "", site.Name, flags), nil
	})
}
