package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/LMBishop/scrapbook/api"
	"github.com/LMBishop/scrapbook/pkg/auth"
	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/constants"
	"github.com/LMBishop/scrapbook/pkg/index"
	"github.com/LMBishop/scrapbook/pkg/server"
	"github.com/LMBishop/scrapbook/web"
)

func main() {
	slog.Info("welcome to scrapbook")

	var cfg config.MainConfig
	cfgStr, err := os.ReadFile(filepath.Join(constants.SysConfDir, "config"))
	if err != nil {
		panic(fmt.Errorf("main config read failed: %w", err))
	}
	cfg, err = config.MainConfigFromString(string(cfgStr))
	if err != nil {
		panic(fmt.Errorf("main config parse failed: %w", err))
	}
	err = config.ValidateMainConfig(&cfg)
	if err != nil {
		panic(fmt.Errorf("main config validation failed: %w", err))
	}

	siteIndex := index.NewSiteIndex()

	sitesDirectory := path.Join(constants.SysDataDir, "sites")
	os.MkdirAll(sitesDirectory, 0o755)

	err = index.ScanDirectory(sitesDirectory, siteIndex)
	if err != nil {
		panic(fmt.Errorf("could not scan data directory: %w", err))
	}

	slog.Info("initial data directory scan complete", "sites", len(siteIndex.GetSites()))

	authenticator := auth.NewAuthenticator()

	mux := http.NewServeMux()

	if cfg.Control.Host == "" {
		slog.Warn("control interface host is empty - neither api or web interface will be accessible")
	} else {
		mux.Handle(fmt.Sprintf("%s/api/", cfg.Control.Host), http.StripPrefix("/api", api.NewMux(&cfg, siteIndex)))
		mux.Handle(fmt.Sprintf("%s/", cfg.Control.Host), web.NewMux(&cfg, siteIndex, authenticator))
	}
	mux.HandleFunc("/", server.ServeSite(siteIndex))

	if cfg.Control.Secret == "" {
		slog.Warn("control interface secret is empty - neither api or web interface will be accessible")
	}

	err = http.ListenAndServe(cfg.Listen, mux)
	slog.Error("http server closing", "reason", err.Error())
}
