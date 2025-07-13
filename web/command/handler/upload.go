package handler

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/upload"
	"github.com/LMBishop/scrapbook/web/command/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetUpload(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}
		site := index.GetSite(siteName)
		if site == nil {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		return html.UploadPage("", "", siteName), nil
	})
}

func PostUpload(mainConfig *config.MainConfig, index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		siteName := r.PathValue("site")
		if siteName == "" {
			return html.ErrorPage("Unknown site: " + siteName), nil
		}

		reader, err := r.MultipartReader()
		if err != nil {
			return html.UploadPage("", fmt.Errorf("Unexpected error: could not read stream: %w", err).Error(), siteName), nil
		}

		version, err := upload.HandleUpload(siteName, reader, index)
		if err != nil {
			return html.UploadPage("", fmt.Errorf("Unexpected error: %w", err).Error(), siteName), nil
		}

		return html.UploadPage(fmt.Sprintf("Version %s created successfully", version), "", siteName), nil
	})
}
