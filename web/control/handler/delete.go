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

func GetDelete(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)
		return html.DeletePage("", "", site.Name), nil
	})
}

func PostDelete(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)

		err := r.ParseForm()
		if err != nil {
			return html.DeletePage("", fmt.Errorf("Could not parse form: %w", err).Error(), site.Name), nil
		}

		if r.FormValue("delete") != "on" {
			return html.DeletePage("", "You need to check the box to continue", site.Name), nil
		}

		err = site.DeleteDataOnDisk()
		if err != nil {
			return html.DeletePage("", fmt.Errorf("Error occurred during data deletion: %w", err).Error(), site.Name), nil
		}
		index.RemoveSite(site.Name)

		return html.DeletePage(fmt.Sprintf("Successfully deleted site %s", site.Name), "", site.Name), nil
	})
}
