package bootstrap

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
)

func newUserConfig(loader *config.Loader) *config.Config {
	cfg := loader.Current()
	if cfg == nil {
		println("[bootstrap] loader.Current() is nil!")
	} else {
		println("[bootstrap] loader.Current() loaded, rules count:", len(cfg.Rules))
	}

	return cfg
}
