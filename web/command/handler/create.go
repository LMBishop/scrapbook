package handler

import (
	"fmt"
	"net/http"
	"path"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/constants"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	"github.com/LMBishop/scrapbook/web/command/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetCreate() func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		return html.CreatePage("", "", html.CreatePageForm{}), nil
	})
}

func PostCreate(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		err := r.ParseForm()
		if err != nil {
			return html.CreatePage("", err.Error(), html.CreatePageForm{}), nil
		}

		name := r.Form.Get("name")
		host := r.Form.Get("host")

		formValues := html.CreatePageForm{Name: name, Host: host}

		if name == "" {
			return html.CreatePage("", "A name must be specified", formValues), nil
		}

		if host == "" {
			return html.CreatePage("", "A host must be specified", formValues), nil
		}

		site, err := site.CreateNewSite(name, path.Join(constants.SysDataDir, "sites"), host)

		if err != nil {
			return html.CreatePage("", fmt.Errorf("Unexpected error: %w", err).Error(), formValues), nil
		}

		index.AddSite(site)

		return html.CreatePage("Successfully created new site", "", formValues), nil
	})
}
