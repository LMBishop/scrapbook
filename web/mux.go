package web

import (
	"net/http"

	"github.com/LMBishop/scrapbook/pkg/auth"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/index"
	h "github.com/LMBishop/scrapbook/web/control/handler"
	m "github.com/LMBishop/scrapbook/web/control/middleware"
)

func NewMux(cfg *config.MainConfig, siteIndex *index.SiteIndex, authenticator *auth.Authenticator) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /authenticate", h.GetAuthenticate())
	mux.HandleFunc("POST /authenticate", h.PostAuthenticate(cfg, authenticator))

	mux.HandleFunc("GET /", m.MustAuthenticate(authenticator, h.GetHome(cfg, siteIndex)))
	mux.HandleFunc("GET /create", m.MustAuthenticate(authenticator, h.GetCreate()))
	mux.HandleFunc("POST /create", m.MustAuthenticate(authenticator, h.PostCreate(cfg, siteIndex)))
	mux.HandleFunc("GET /site/{site}/", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.GetSite(cfg, siteIndex))))
	mux.HandleFunc("GET /site/{site}/upload", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.GetUpload(siteIndex))))
	mux.HandleFunc("POST /site/{site}/upload", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.PostUpload(cfg, siteIndex))))
	mux.HandleFunc("GET /site/{site}/flags", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.GetFlags(siteIndex))))
	mux.HandleFunc("POST /site/{site}/flags", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.PostFlags(cfg, siteIndex))))
	mux.HandleFunc("GET /site/{site}/config", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.GetConfig(siteIndex))))
	mux.HandleFunc("POST /site/{site}/config", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.PostConfig(cfg, siteIndex))))
	mux.HandleFunc("GET /site/{site}/delete", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.GetDelete(siteIndex))))
	mux.HandleFunc("POST /site/{site}/delete", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, h.PostDelete(cfg, siteIndex))))

	mux.HandleFunc("GET /site/{site}/version/{version}/", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, m.WithVersion(h.GetVersion(siteIndex)))))
	mux.HandleFunc("GET /site/{site}/version/{version}/delete", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, m.WithVersion(h.GetVersionDelete(siteIndex)))))
	mux.HandleFunc("POST /site/{site}/version/{version}/delete", m.MustAuthenticate(authenticator, m.WithSite(siteIndex, m.WithVersion(h.PostVersionDelete(cfg, siteIndex)))))

	return mux
}
