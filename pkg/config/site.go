package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type SiteFlag uint

const (
	FlagDisable SiteFlag = 1 << iota
	FlagTLS
	FlagIndex
	FlagPassword
	FlagReadOnly
)

type SiteConfig struct {
	Host string

	Flags SiteFlag
}

func ReadSiteConfig(filePath string, dst *SiteConfig) error {
	config, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(config, dst); err != nil {
		return err
	}
	return nil
}

func WriteSiteConfig(filePath string, src *SiteConfig) error {
	config, err := toml.Marshal(src)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, config, 0o644); err != nil {
		return err
	}
	return nil
}
