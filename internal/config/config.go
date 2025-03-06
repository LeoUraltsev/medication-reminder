package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"
)

type Config struct {
	App      App      `yaml:"app"`
	Http     Http     `yaml:"http"`
	Postgres Postgres `yaml:"postgres"`
}

type App struct {
	LogLevel string `yaml:"loglevel" env-default:"debug"`
}

type Http struct {
	Address      string        `yaml:"address" env-default:"localhost:8008"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"5s"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"5s"`
}

type Postgres struct {
	ConnectionURL string `yaml:"connection_url" env-required:"true" env:"MR_PG_CONNECTION_URL"`
}

func init() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Info("Not found .env file")
	}
}

func New() (*Config, error) {
	var cfg Config
	var cfgPath string
	flag.StringVar(&cfgPath, "config-path", "", "path to app config file")
	flag.Parse()

	if cfgPath == "" {
		cfgPath = os.Getenv("MR_CONFIG_PATH")

	}

	if _, err := os.Stat(cfgPath); err != nil {
		return nil, fmt.Errorf("config file not found: %s", cfgPath)
	}

	if err := cleanenv.ReadConfig(cfgPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed parse config file: %w", err)
	}

	return &cfg, nil
}
