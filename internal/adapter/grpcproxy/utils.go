package grpcproxy

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

// parsePath: /full.Service/Method -> service, method
func parsePath(p string) (string, string) {
	if !strings.HasPrefix(p, "/") {
		return "", ""
	}

	parts := strings.SplitN(p[1:], "/", 2)
	if len(parts) != 2 {
		return "", ""
	}

	return parts[0], parts[1]
}

// metadataFromHeader converts HTTP headers to gRPC metadata.MD
func metadataFromHeader(h http.Header) metadata.MD {
	md := metadata.MD{}

	for k, vals := range h {
		lk := strings.ToLower(k)
		md[lk] = append(md[lk], vals...)
	}

	return md
}

// isGRPCRequest checks if the request is a gRPC request based on the content type
func writeHTTPError(w http.ResponseWriter, statusCode int, msg string) {
	http.Error(w, msg, statusCode)
}

// isReflectionRequest checks if the request is for gRPC reflection service
func isReflectionRequest(svc, mth string) bool {
	return (svc == "grpc.reflection.v1.ServerReflection" || svc == "grpc.reflection.v1alpha.ServerReflection") && mth == "ServerReflectionInfo"
}

// isGRPCRequest checks if the request is a gRPC request
func isGRPCRequest(r *http.Request) bool {
	return r.Header.Get("content-type") == "application/grpc"
}

// dialBackend establishes a gRPC connection to the backend service
func dialBackend(addr string, log *zap.Logger) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Warn("dial backend", zap.Error(err))

		return nil, err
	}

	return conn, nil
}

// errMsg returns the error message if err is not nil, otherwise returns the fallback message
func errMsg(err error, fallback string) string {
	if err != nil {
		return err.Error()
	}

	return fallback
}

// writeGRPCOK writes a successful gRPC response to the HTTP response writer
func writeGRPCOK(w http.ResponseWriter, payload []byte) {
	w.Header().Set("content-type", "application/grpc")
	w.Header().Set("grpc-status", "0")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

// writeGRPCError writes a gRPC error response to the HTTP response writer
func writeGRPCError(w http.ResponseWriter, code codes.Code, msg string) {
	w.Header().Set("content-type", "application/grpc")
	w.Header().Set("grpc-status", strconv.Itoa(int(code)))
	w.Header().Set("grpc-message", msg)
	w.WriteHeader(http.StatusOK)
}

// writeStreamRecvError handles errors that occur during receiving messages from a gRPC stream
func writeStreamRecvError(w http.ResponseWriter, err error) {
	if st, ok := status.FromError(err); ok {
		writeGRPCError(w, st.Code(), st.Message())
	} else {
		writeGRPCError(w, codes.Internal, err.Error())
	}
}
