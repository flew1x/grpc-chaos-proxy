package grpcproxy

import (
	"context"
	"errors"
	"github.com/flew1x/grpc-chaos-proxy/internal/apperr"
	"github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"github.com/flew1x/grpc-chaos-proxy/internal/entity"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
)

type Server struct {
	engine *engine.Engine
	cfg    *config.Config
	log    *zap.Logger
}

func New(e *engine.Engine, c *config.Config, l *zap.Logger) *Server {
	if e == nil || c == nil || l == nil {
		panic("engine, config, logger must not be nil")
	}

	return &Server{engine: e, cfg: c, log: l}
}

func (s *Server) Run(lc fx.Lifecycle) {
	h2s := &http2.Server{}
	handler := http.HandlerFunc(s.handle)

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			lis, err := net.Listen("tcp", s.cfg.Listener.Address)
			if err != nil {
				return err
			}

			addr := lis.Addr().String()
			s.log.Info("proxy listening", zap.String("addr", addr))

			go func() {
				s.log.Info("proxy serve goroutine started", zap.String("addr", addr))
				if err := http.Serve(lis, h2c.NewHandler(handler, h2s)); err != nil {
					s.log.Error("proxy server stopped", zap.Error(err))
				}
			}()

			s.log.Info("OnStart completed, server should be running", zap.String("addr", addr))

			return nil
		},
		OnStop: func(ctx context.Context) error {
			s.log.Info("proxy OnStop called")

			return nil
		},
	})
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	logger := s.log.With(
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.Any("headers", r.Header),
	)

	svc, mth := parsePath(r.URL.Path)
	if svc == "" || mth == "" {
		logger.Error("bad grpc path or method missing")
		writeHTTPError(w, http.StatusBadRequest, "bad grpc path")

		return
	}

	if isReflectionRequest(svc, mth) {
		logger.Error("reflection/streaming not supported", zap.String("service", svc))
		writeHTTPError(w, http.StatusNotImplemented, "gRPC reflection and streaming not supported by proxy")

		return
	}

	frame := &engine.Frame{
		Ctx:       r.Context(),
		Service:   svc,
		Method:    mth,
		MD:        metadataFromHeader(r.Header),
		Direction: entity.DirectionInbound,
	}

	if err := s.engine.Process(frame); err != nil {
		if errors.Is(err, apperr.ErrNoMatchingRule) {
			writeGRPCError(w, codes.NotFound, "no matching rule for this method")
			logger.Warn("no matching rule for request", zap.String("service", svc))

			return
		}

		if st, ok := status.FromError(err); ok {
			writeGRPCError(w, st.Code(), st.Message())
			logger.Error("engine.Process gRPC error", zap.Error(err), zap.String("service", svc))

			return
		}

		writeGRPCError(w, codes.Internal, "internal chaos error")
		logger.Error("engine.Process unknown error", zap.Error(err), zap.String("service", svc))

		return
	}

	logger.Info("proxying request", zap.String("service", svc))

	Proxy(w, r, frame, logger, s.cfg)
}
