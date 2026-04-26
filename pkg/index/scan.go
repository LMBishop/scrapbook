package index

import (
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
		site := &site.Site{
			Name: siteName,
			Path: sitePath,
		}

		cfgStr, err := site.ReadSiteConfigFile()
		if err != nil {
			slog.Warn("failed to read site config", "site", siteName, "reason", err)
		} else {
			cfg, err := config.SiteConfigFromString(cfgStr)
			if err != nil {
				slog.Warn("failed to parse site config", "site", siteName, "reason", err)
			} else {
				site.Config = &cfg
			}
		}

		site.Flags = site.ReadSiteFlags()
		site.Initialise()
		dst.AddSite(site)
	}

	return nil
}
