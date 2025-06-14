package bootstrap

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"go.uber.org/zap"
)

func newEngine(loader *config.Loader, logger *zap.Logger) (*engine.Engine, error) {
	return engine.New(loader, logger)
}
