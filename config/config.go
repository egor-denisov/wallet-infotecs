package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	// Config -.
	Config struct {
		App        `yaml:"app"`
		HTTP       `yaml:"http"`
		Log        `yaml:"logger"`
		PG         `yaml:"postgres"`
	}

	// App -.
	App struct {
		Name           string  `env-required:"true" yaml:"name"            env:"APP_NAME"`
		Version        string  `env-required:"true" yaml:"version"         env:"APP_VERSION"`
		DefaultBalance float64 `env-required:"true" yaml:"default_balance" env:"DEFAULT_BALANCE"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		User     string `env-required:"true"  yaml:"pg_user"  env:"POSTGRES_USER"`
		DB       string `env-required:"true" yaml:"pg_db" env:"POSTGRES_DB"`
		Host     string `env-required:"true" yaml:"pg_host" env:"POSTGRES_HOST"`
		Port     int    `env-required:"true" yaml:"pg_port" env:"POSTGRES_PORT"`
		Password string `env-required:"true" yaml:"pg_password" env:"POSTGRES_PASSWORD"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := godotenv.Load()
    if err != nil {
        fmt.Println("Error loading .env file")
    }

	err = cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
