package initiator

import (
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/abort"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/chaos"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/code"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/disconnect"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/header"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/network"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/ratelimit"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/script"
	_ "github.com/flew1x/grpc-chaos-proxy/internal/core/injector/spammer"
)
