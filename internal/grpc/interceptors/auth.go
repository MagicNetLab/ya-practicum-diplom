package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
	}

	var token string
	values := md.Get("token")
	if len(values) == 1 {
		token = values[0]
	}

	cnf := config.GetJWTConfig()
	if token != "" && jwt.VerifyToken(token, cnf.GetJWTSecret()) {
		return handler(ctx, req)
	}

	return nil, status.Errorf(codes.Unauthenticated, "unauthorized")
}

func GuestInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return handler(ctx, req)
	}

	var token string

	values := md.Get("token")
	if len(values) == 1 {
		token = values[0]
	}

	cnf := config.GetJWTConfig()
	if token != "" && jwt.VerifyToken(token, cnf.GetJWTSecret()) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden for authorized users")
	}

	return handler(ctx, req)
}
