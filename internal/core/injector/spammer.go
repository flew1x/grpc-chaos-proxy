package injector

import (
	"context"
	"fmt"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
	"time"
)

type SpammerInjector struct {
	count       int
	delay       *config.DelayAction
	backendAddr string
	sender      func(ctx context.Context, service, method string, md metadata.MD) error
}

var (
	proxyAddrOnce   sync.Once
	cachedProxyAddr string
)

func getProxyAddr() string {
	proxyAddrOnce.Do(func() {
		cfg := config.GetCurrentConfig()
		if cfg != nil && cfg.Listener.Address != "" {
			cachedProxyAddr = cfg.Listener.Address
		}
	})

	return cachedProxyAddr
}

func NewSpammerInjector(cfg any) (engine.Injector, error) {
	sa, ok := cfg.(*config.SpammerAction)
	if !ok || sa == nil || sa.Count <= 0 {
		return nil, fmt.Errorf("invalid config for SpammerInjector")
	}

	proxyAddr := getProxyAddr()

	return &SpammerInjector{
		count:       sa.Count,
		delay:       sa.DelayAction,
		backendAddr: proxyAddr,
		sender:      nil,
	}, nil
}

func (s *SpammerInjector) Apply(f *engine.Frame) error {
	if f.MD != nil {
		if vals := f.MD.Get("x-spammer-request"); len(vals) > 0 && vals[0] == "1" {
			return nil
		}
	}

	send := s.sender

	if send == nil {
		send = func(ctx context.Context, service, method string, md metadata.MD) error {
			conn, err := grpc.NewClient(s.backendAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				return err
			}

			defer conn.Close()

			fullMethod := "/" + service + "/" + method
			ctx = metadata.NewOutgoingContext(ctx, md)

			return conn.Invoke(ctx, fullMethod, &emptypb.Empty{}, &emptypb.Empty{})
		}
	}

	for i := 0; i < s.count; i++ {
		if s.delay != nil && (s.delay.MinMS > 0 || s.delay.MaxMS > 0) {
			minVal := s.delay.MinMS
			maxVal := s.delay.MaxMS

			if maxVal < minVal {
				minVal, maxVal = maxVal, minVal
			}

			dur := time.Duration(minVal)
			if maxVal > minVal {
				dur = time.Duration(minVal + (randInt(maxVal - minVal)))
			}

			time.Sleep(dur * time.Millisecond)
		}

		md := metadata.New(map[string]string{"x-spammer-request": "1"})
		if len(f.MD) > 0 {
			for k, v := range f.MD {
				if len(v) > 0 {
					md.Set(k, v[0])
				}
			}
		}

		_ = send(f.Ctx, f.Service, f.Method, md)
	}
	return nil
}

func randInt(n int) int {
	if n <= 0 {
		return 0
	}

	return int(time.Now().UnixNano() % int64(n))
}

func init() {
	engine.Register(entity.SpammerType, NewSpammerInjector)
}
