package config

import (
	"fmt"
	"strings"

	"codeberg.org/emersion/go-scfg"
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
	Host      string `scfg:"host"`
	Password  string `scfg:"password"`
	Retention struct {
		Amount uint `scfg:"amount"`
	} `scfg:"retention"`
}

const DefaultSiteConfig = `host "%s"
`

func NewSiteConfig(host string) (SiteConfig, string) {
	var config SiteConfig
	strConfig := fmt.Sprintf(DefaultSiteConfig, host)
	scfg.NewDecoder(strings.NewReader(strConfig)).Decode(&config)
	return config, strConfig
}

func SiteConfigFromString(str string) (SiteConfig, error) {
	var config SiteConfig
	err := scfg.NewDecoder(strings.NewReader(str)).Decode(&config)
	return config, err
}
