package api

import (
	"net/http"

	"github.com/LMBishop/scrapbook/api/handler"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
)

func NewMux(cfg *config.MainConfig, siteIndex *index.SiteIndex) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /site/{site}/upload", handler.UploadSiteVersion(cfg, siteIndex))

	return mux
}
