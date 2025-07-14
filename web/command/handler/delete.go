package handler

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/command/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetDelete(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		return html.DeletePage("", "", siteName), nil
	})
}

func PostDelete(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
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
			return html.DeletePage("", fmt.Errorf("Could not parse form: %w", err).Error(), siteName), nil
		}

		if r.FormValue("delete") != "on" {
			return html.DeletePage("", "You need to check the box to continue", siteName), nil
		}

		err = site.DeleteDataOnDisk()
		if err != nil {
			return html.DeletePage("", fmt.Errorf("Error occurred during data deletion: %w", err).Error(), siteName), nil
		}
		index.RemoveSite(siteName)

		return html.DeletePage(fmt.Sprintf("Successfully deleted site %s", siteName), "", siteName), nil
	})
}
