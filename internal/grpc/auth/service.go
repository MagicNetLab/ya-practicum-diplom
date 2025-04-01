package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/encryptor"
)

// MakeService возвращает настроенный сервис авторизации
func MakeService(cnf config.AppConfig, repo repository.AuthRepository) (Service, error) {
	return Service{cnf: cnf, store: repo}, nil
}

// Service  сервис авторизации
type Service struct {
	pb.AuthServer
	cnf   config.AppConfig
	store repository.AuthRepository
}

// Auth авторизация пользователя
func (s *Service) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, err := s.store.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	// Verify password
	hashedPassword, err := encryptor.EncryptPassword(req.Secret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to process password")
	}

	if hashedPassword != user.GetPassword() {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid credentials")
	}

	token, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = s.store.CreateToken(ctx, user.GetUID(), token, false, time.Now().Add(s.cnf.JWTTokenLifeTime()))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create token")
	}

	return &pb.AuthResponse{Token: token, RefreshToken: ""}, nil
}

// Register регистрация пользователя
func (s *Service) Register(ctx context.Context, req *pb.RegRequest) (*pb.RegResponse, error) {
	hasLogin, err := s.store.HasLogin(ctx, req.Login)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if hasLogin {
		return nil, status.Errorf(codes.AlreadyExists, "Login already exists")
	}

	user, err := s.store.CreateUser(ctx, req.Login, req.Secret)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	token, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshToken, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.RegResponse{Token: token, RefreshToken: refreshToken}, nil
}

// Refresh обновление токена
func (s *Service) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	if isValidToken := jwt.VerifyToken(req.Token, s.cnf.JWTSecret()); !isValidToken {
		return nil, status.Errorf(codes.PermissionDenied, "Invalid token1")
	}

	tokenData, err := jwt.ParseToken(req.Token, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	checkToken, err := s.store.HasToken(ctx, req.Token)
	if err != nil || !checkToken {
		return nil, status.Errorf(codes.Internal, "Invalid token")
	}

	user, err := s.store.GetUserByUID(ctx, tokenData.UID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	token, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshToken, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &pb.RefreshResponse{Token: token, RefreshToken: refreshToken}, nil
}

// RegisterService регистрирует сервис в gRPC сервере
func RegisterService(gRPCServer *grpc.Server, s Service) {
	pb.RegisterAuthServer(gRPCServer, &s)
}
