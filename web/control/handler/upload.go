package handler

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	"github.com/LMBishop/scrapbook/pkg/upload"
	"github.com/LMBishop/scrapbook/web/control/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetUpload(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)
		return html.UploadPage("", "", site.Name), nil
	})
}

func PostUpload(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		site := r.Context().Value("site").(*site.Site)

		reader, err := r.MultipartReader()
		if err != nil {
			return html.UploadPage("", fmt.Errorf("Unexpected error: could not read stream: %w", err).Error(), site.Name), nil
		}

		version, err := upload.HandleUpload(site.Name, "WebUI", reader, index)
		if err != nil {
			return html.UploadPage("", fmt.Errorf("Unexpected error: %w", err).Error(), site.Name), nil
		}

		return html.UploadPage(fmt.Sprintf("Version %s created successfully", version), "", site.Name), nil
	})
}
