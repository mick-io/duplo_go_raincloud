package config

import (
	"path/filepath"
	"strings"

	"github.com/go-playground/validator"
	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string `validate:"required"`
	Name     string `validate:"required"`
	Password string `validate:"required"`
	Port     int    `validate:"required,min=1024,max=65535"`
	User     string `validate:"required"`
}

type Config struct {
	Database *DatabaseConfig
	Server   struct {
		Environment string `validate:"required"`
		Port        int    `validate:"required,min=1024,max=65535"`
	}
	API struct {
		ForecastAPIBaseURL string `mapstructure:"forecast_api_base_url" validate:"required,url"`
	}
}

func (c *Config) validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func Load(path string) (*Config, error) {
	var config Config

	dirpath, filename := filepath.Split(path)
	ext := filepath.Ext(filename)

	viper.AddConfigPath(dirpath)
	viper.SetConfigName(strings.TrimSuffix(filename, ext))
	viper.SetConfigType(strings.TrimPrefix(ext, "."))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return &config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return &config, err
	}

	if err := config.validate(); err != nil {
		return &config, err
	}

	return &config, nil
}
