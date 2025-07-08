package index

import (
	"fmt"
	"log/slog"
	"os"
	"path"

	"github.com/LMBishop/scrapbook/pkg/config"
	"github.com/LMBishop/scrapbook/pkg/site"
)

func ScanDirectory(dir string, dst *SiteIndex) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}

		siteName := e.Name()
		sitePath := path.Join(dir, siteName)
		cfg, err := readSiteConfig(sitePath)
		if err != nil {
			slog.Warn("failed to read site", "site", siteName, "reason", err)
			continue
		}

		site := site.NewSite(siteName, sitePath, cfg)
		dst.AddSite(site)
	}

	return nil
}

func readSiteConfig(dir string) (*config.SiteConfig, error) {
	siteFile := path.Join(dir, "site.toml")
	cfg := &config.SiteConfig{}
	err := config.ReadSiteConfig(siteFile, cfg)
	if err != nil {
		return nil, fmt.Errorf("site file invalid: %s", err)
	}
	return cfg, nil
}
