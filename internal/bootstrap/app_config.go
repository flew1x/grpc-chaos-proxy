package bootstrap

import (
	"github.com/joho/godotenv"
	"os"
)

type AppConfig struct {
	ConfigPath string
	LogLevel   string
}

func loadAppConfig() (*AppConfig, error) {
	_ = godotenv.Load()
	cfg := &AppConfig{
		ConfigPath: os.Getenv("CONFIG_PATH"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}

	if cfg.ConfigPath == "" {
		return nil, &os.PathError{
			Op: "get CONFIG_PATH",
		}
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}
