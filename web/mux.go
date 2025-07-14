package web

import (
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/auth"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/web/command/handler"
	"github.com/LMBishop/scrapbook/web/command/middleware"
)

func NewMux(cfg *config.MainConfig, siteIndex *index.SiteIndex, authenticator *auth.Authenticator) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /authenticate", handler.GetAuthenticate())
	mux.HandleFunc("POST /authenticate", handler.PostAuthenticate(cfg, authenticator))

	mux.HandleFunc("GET /", middleware.MustAuthenticate(authenticator, handler.GetHome(cfg, siteIndex)))
	mux.HandleFunc("GET /create", middleware.MustAuthenticate(authenticator, handler.GetCreate()))
	mux.HandleFunc("POST /create", middleware.MustAuthenticate(authenticator, handler.PostCreate(cfg, siteIndex)))
	mux.HandleFunc("GET /site/{site}/", middleware.MustAuthenticate(authenticator, handler.GetSite(cfg, siteIndex)))
	mux.HandleFunc("GET /site/{site}/upload", middleware.MustAuthenticate(authenticator, handler.GetUpload(siteIndex)))
	mux.HandleFunc("POST /site/{site}/upload", middleware.MustAuthenticate(authenticator, handler.PostUpload(cfg, siteIndex)))
	mux.HandleFunc("GET /site/{site}/flags", middleware.MustAuthenticate(authenticator, handler.GetFlags(siteIndex)))
	mux.HandleFunc("POST /site/{site}/flags", middleware.MustAuthenticate(authenticator, handler.PostFlags(cfg, siteIndex)))
	mux.HandleFunc("GET /site/{site}/host", middleware.MustAuthenticate(authenticator, handler.GetHost(siteIndex)))
	mux.HandleFunc("POST /site/{site}/host", middleware.MustAuthenticate(authenticator, handler.PostHost(cfg, siteIndex)))

	return mux
}
