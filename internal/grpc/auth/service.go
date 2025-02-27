package auth

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/conf"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"

	"google.golang.org/grpc/codes"
	_ "google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repo"
)

// Service  сервис авторизации
type Service struct {
	pb.AuthServer
	cnf   conf.Configurator
	store repo.Repository
}

// Auth авторизация пользователя
func (s *Service) Auth(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	user, err := s.store.GetUserByLoginAndPassword(ctx, req.Login, req.Secret)
	if err != nil {
		return nil, status.Errorf(codes.PermissionDenied, err.Error())
	}

	token, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	refreshToken, err := jwt.GenerateToken(user, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	err = s.store.CreateToken(ctx, token, user.GetUID(), false)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create token")
	}

	err = s.store.CreateToken(ctx, refreshToken, user.GetUID(), true)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create refresh token")
	}

	return &pb.AuthResponse{Token: token, RefreshToken: refreshToken}, nil
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

	checkToken, err := s.store.HasToken(ctx, req.Token, tokenData.UID, true)
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
