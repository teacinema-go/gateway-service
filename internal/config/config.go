package config

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"github.com/teacinema-go/core/constants"
)

type Config struct {
	App     AppConfig     `mapstructure:",squash"`
	Service ServiceConfig `mapstructure:",squash"`
}

type AppConfig struct {
	Env  constants.Env `mapstructure:"APP_ENV" validate:"required"`
	Port int           `mapstructure:"APP_PORT" validate:"required"`
	Host string        `mapstructure:"APP_HOST" validate:"required,url"`
}

type ServiceConfig struct {
	AuthServiceURL string `mapstructure:"SERVICE_AUTH_URL" validate:"required,url"`
}

func Load() (*Config, error) {
	viper.SetDefault("APP_ENV", constants.Development)
	viper.SetDefault("APP_PORT", 8000)
	viper.SetDefault("APP_HOST", "http://localhost:8000")
	viper.SetDefault("SERVICE_AUTH_URL", "http://localhost:50051")

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &cfg, nil
}
