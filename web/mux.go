package web

import (
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/command/handler"
)

func NewMux(cfg *config.MainConfig, siteIndex *index.SiteIndex) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.GetHome(cfg, siteIndex))
	mux.HandleFunc("GET /create", handler.GetCreate())
	mux.HandleFunc("POST /create", handler.PostCreate(cfg, siteIndex))
	mux.HandleFunc("GET /site/{site}/", handler.GetSite(cfg, siteIndex))
	mux.HandleFunc("GET /site/{site}/upload", handler.GetUpload(siteIndex))
	mux.HandleFunc("POST /site/{site}/upload", handler.PostUpload(cfg, siteIndex))

	return mux
}
