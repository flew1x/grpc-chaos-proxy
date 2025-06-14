package bootstrap

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"go.uber.org/zap"
)

func newLoaderConfig(cfg *AppConfig, logger *zap.Logger) (*config.Loader, error) {
	if cfg.ConfigPath == "" {
		logger.Fatal("CONFIG_PATH environment variable not set")
	}

	loader, err := config.NewLoader(cfg.ConfigPath, logger)
	if err != nil {
		return nil, err
	}

	return loader, nil
}
