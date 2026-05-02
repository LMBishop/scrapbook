package handler

import (
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/control/html"
	. "maragu.dev/gomponents"
	ghttp "maragu.dev/gomponents/http"
)

func GetVersion(index *index.SiteIndex) func(http.ResponseWriter, *http.Request) {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		versionMetadata := r.Context().Value("versionMetadata").(config.VersionMeta)
		return html.VersionPage(versionMetadata), nil
	})
}
