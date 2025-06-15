package utils

import (
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"sync"
	"time"
)

var (
	proxyAddrOnce   sync.Once
	cachedProxyAddr string
)

func GetProxyAddr() string {
	proxyAddrOnce.Do(func() {
		cfg := config.GetCurrentConfig()
		if cfg != nil && cfg.Listener.Address != "" {
			cachedProxyAddr = cfg.Listener.Address
		}
	})

	return cachedProxyAddr
}

func RandInt(n int) int {
	if n <= 0 {
		return 0
	}

	return int(time.Now().UnixNano() % int64(n))
}
