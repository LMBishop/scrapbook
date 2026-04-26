package server

import (
	"github.com/LMBishop/scrapbook/pkg/html"
	"github.com/LMBishop/scrapbook/pkg/index"
	"net/http"
)

func ServeSite(siteIndex *index.SiteIndex) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		site := siteIndex.GetSiteByHost(r.Host)
		if site == nil || site.Handler() == nil {
			w.WriteHeader(http.StatusNotFound)
			html.NotFoundSitePage(r.Host).Render(w)
			return
		}

		site.Handler().ServeHTTP(w, r)
	}
}
