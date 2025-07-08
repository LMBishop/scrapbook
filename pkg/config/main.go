package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/pelletier/go-toml/v2"
)

type MainConfig struct {
	Listen struct {
		Address string `validate:"required,ip"`
		Port    uint16 `validate:"required"`
	}

	Command struct {
		Host   string
		Secret string

		API struct {
			Enable bool
		}

		Web struct {
			Enable bool
		}
	}
}

func ReadMainConfig(filePath string, dst *MainConfig) error {
	config, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(config, dst); err != nil {
		return err
	}
	return nil
}

func ValidateMainConfig(cfg *MainConfig) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(cfg)
}
