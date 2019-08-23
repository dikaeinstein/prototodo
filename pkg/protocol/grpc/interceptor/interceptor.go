package interceptor

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
)

// LogRPCCalls logs all RPC methods call.
func LogRPCCalls(l *zap.Logger) grpc.UnaryServerInterceptor {
	return grpc_zap.UnaryServerInterceptor(l)
}
