package grpcproxy

import (
	"context"
	"io"
	"net/http"
	"time"

	cfg "github.com/flew1x/grpc-chaos-proxy/internal/config"
	"github.com/flew1x/grpc-chaos-proxy/internal/core/engine"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	defaultGRPCTimeout = 5 * time.Second
)

// Proxy handles ONLY unary-gRPC requests (application/grpc).
// It does not attempt to understand the protobuf schema; instead,
// it wraps the input payload in google.protobuf.Any, which is
// enough for most e2e tests
func Proxy(
	w http.ResponseWriter,
	r *http.Request,
	f *engine.Frame,
	log *zap.Logger,
	c *cfg.Config,
) {
	if !isGRPCRequest(r) {
		writeGRPCError(w, codes.Unimplemented, "only unary gRPC supported")

		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), defaultGRPCTimeout)
	defer cancel()

	conn, err := dialBackend(c.Backend.Address, log)
	if err != nil {
		writeGRPCError(w, codes.Unavailable, "backend unavailable")

		return
	}
	defer conn.Close()

	raw, err := io.ReadAll(r.Body)
	if err != nil || len(raw) == 0 {
		writeGRPCError(w, codes.Internal, "read body: "+errMsg(err, "empty body"))

		return
	}

	in := &anypb.Any{Value: raw}
	fullMethod := "/" + f.Service + "/" + f.Method
	desc := &grpc.StreamDesc{ServerStreams: false, ClientStreams: false}

	cs, err := conn.NewStream(ctx, desc, fullMethod)
	if err != nil {
		writeGRPCError(w, codes.Unavailable, "open stream: "+err.Error())

		return
	}

	if err := cs.SendMsg(in); err != nil {
		writeGRPCError(w, codes.Internal, "send body: "+err.Error())

		return
	}

	_ = cs.CloseSend()

	out := &anypb.Any{}

	if err := cs.RecvMsg(out); err != nil {
		writeStreamRecvError(w, err)

		return
	}

	writeGRPCOK(w, out.Value)
}
