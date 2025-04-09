package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

// LoggerInterceptor логирование всех запросов
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	logger.Info("grpc request info",
		logger.StrArg("method", info.FullMethod),
		logger.TimeDurationArg("duration", duration))

	return resp, err
}
