package middleware

import (
	"context"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/site"
	"github.com/LMBishop/scrapbook/web/control/html"
)

func WithSite(index *index.SiteIndex, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var siteName string
		var site *site.Site

		siteName = r.PathValue("site")
		if siteName == "" {
			goto siteUnknown
		}
		site = index.GetSite(siteName)
		if site == nil {
			goto siteUnknown
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "site", site)))
		return

	siteUnknown:
		html.ErrorPage("Unknown site: " + siteName).Render(w)
	})
}

func WithVersion(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var versionMetadata config.VersionMeta
		var versionMetadataStr string
		var err error

		site := r.Context().Value("site").(*site.Site)

		versionName := r.PathValue("version")
		if versionName == "" {
			goto versionUnknown
		}
		versionMetadataStr, err = site.ReadVersionMetadata(versionName)
		if err != nil {
			goto versionErr
		}
		versionMetadata, err = config.VersionMetaFromString(versionMetadataStr)
		if err != nil {
			goto versionErr
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "versionMetadata", versionMetadata)))
		return

	versionUnknown:
		html.ErrorPage("Unknown version: " + versionName).Render(w)
		return

	versionErr:
		html.ErrorPage("Could not read version metadata: " + err.Error()).Render(w)
	})
}
