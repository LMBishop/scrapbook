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

func GetVersionDelete(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		version := r.Context().Value("versionMetadata").(config.VersionMeta)
		return html.DeletePage("", "", version.Hash[:8]), nil
	})
}

func PostVersionDelete(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)
		version := r.Context().Value("versionMetadata").(config.VersionMeta)

		err := r.ParseForm()
		if err != nil {
			return html.DeletePage("", fmt.Errorf("Could not parse form: %w", err).Error(), version.Hash[:8]), nil
		}

		if r.FormValue("delete") != "on" {
			return html.DeletePage("", "You need to check the box to continue", version.Hash[:8]), nil
		}

		err = site.DeleteVersion(version.Hash)
		if err != nil {
			return html.DeletePage("", fmt.Errorf("Error occurred during data deletion: %w", err).Error(), version.Hash[:8]), nil
		}

		return html.DeletePage(fmt.Sprintf("Successfully deleted site %s", version.Hash[:8]), "", version.Hash[:8]), nil
	})
}
