package bootstrap

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/adapter/grpcproxy"
	"go.uber.org/fx"
)

// Module wires up all core dependencies for the proxy app
func Module() fx.Option {
	return fx.Options(
		fx.Provide(
			loadAppConfig,
			newLogger,
			newLoaderConfig,
			newEngine,
			newUserConfig,
		),
		fx.Provide(grpcproxy.New),
		fx.Invoke((*grpcproxy.Server).Run),
	)
}

func NewApp() *fx.App {
	return fx.New(Module())
}
