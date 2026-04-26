package config

import (
	"strings"

	"codeberg.org/emersion/go-scfg"
	"github.com/go-playground/validator/v10"
)

type MainConfig struct {
	Listen string `scfg:"listen" validate:"required"`

	Control struct {
		Host   string `scfg:"host"`
		Secret string `scfg:"secret"`

		API struct {
			Enable bool `scfg:"enable"`
		} `scfg:"api"`

		Web struct {
			Enable bool `scfg:"enable"`
		} `scfg:"web"`
	} `scfg:"control"`
}

func MainConfigFromString(str string) (MainConfig, error) {
	var config MainConfig
	err := scfg.NewDecoder(strings.NewReader(str)).Decode(&config)
	return config, err
}

func ValidateMainConfig(cfg *MainConfig) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}
