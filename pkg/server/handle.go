package server

import (
	"fmt"
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/index"
)

func ServeSite(siteIndex *index.SiteIndex) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		site := siteIndex.GetSiteByHost(r.Host)
		if site == nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "unknown host %s", r.Host)
			return
		}

		site.Handler.ServeHTTP(w, r)
	}
}
