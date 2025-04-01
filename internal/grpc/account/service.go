package account

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// MakeService возвращает настроенный сервис аккаунтов
func MakeService(repo repository.AccountRepository, cnf config.AppConfig) (Service, error) {
	return Service{store: repo, cnf: cnf}, nil
}

// Service  сервис аккаунтов
type Service struct {
	pb.AccountsServer
	store repository.AccountRepository
	cnf   config.AppConfig
}

// Get получение аккаунта по идентификатору
func (s *Service) Get(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	userUID, err := jwt.GetUIDFromContext(ctx, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err = uuid.Validate(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	result, err := s.store.GetAccount(ctx, req.GetId(), userUID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "account not found")
	}

	acc := &pb.Account{
		Id:          result.GetID(),
		Uid:         result.GetUID(),
		Login:       result.GetLogin(),
		Password:    result.GetPassword(),
		Url:         result.GetURL(),
		Description: result.GetDescription(),
	}

	return &pb.GetAccountResponse{Account: acc}, nil

}

// Create создание аккаунта
func (s *Service) Create(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	userUID, err := jwt.GetUIDFromContext(ctx, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	acc, err := s.store.CreateAccount(ctx, userUID, req.GetLogin(), req.GetPassword(), req.GetUrl(), req.GetDescription())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "account not created")
	}

	account := &pb.Account{
		Id:          acc.GetID(),
		Uid:         acc.GetUID(),
		Login:       acc.GetLogin(),
		Password:    acc.GetPassword(),
		Url:         acc.GetURL(),
		Description: acc.GetDescription(),
	}

	return &pb.CreateAccountResponse{Acc: account}, nil
}

// Remove удаление аккаунта по идентификатору
func (s *Service) Remove(ctx context.Context, req *pb.RemoveAccountRequest) (*pb.RemoveAccountResponse, error) {
	userUID, err := jwt.GetUIDFromContext(ctx, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	id := req.GetId()
	err = s.store.RemoveAccount(ctx, id, userUID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "account not removed")
	}

	return &pb.RemoveAccountResponse{}, nil
}

// Search поиск аккаунтов по критериям
func (s *Service) Search(ctx context.Context, req *pb.SearchAccountRequest) (*pb.SearchAccountResponse, error) {
	userUID, err := jwt.GetUIDFromContext(ctx, s.cnf.JWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := models.AccountSearch{}
	search.UID = userUID

	if req.GetSearch() != "" {
		search.Search = req.GetSearch()
	}

	if req.GetLimit() != "" {
		limit, err := strconv.Atoi(req.GetLimit())
		if err == nil {
			search.Limit = limit
		}
	}

	if req.GetOffset() != "" {
		offset, err := strconv.Atoi(req.GetOffset())
		if err == nil {
			search.Offset = offset
		}
	}

	results, err := s.store.SearchAccounts(ctx, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "account not removed")
	}
	var accounts []*pb.Account
	for _, a := range results {
		accounts = append(accounts, &pb.Account{
			Id:          a.GetID(),
			Uid:         a.GetUID(),
			Login:       a.GetLogin(),
			Password:    a.GetPassword(),
			Url:         a.GetURL(),
			Description: a.GetDescription(),
		})
	}

	return &pb.SearchAccountResponse{Acc: accounts}, nil
}

// RegisterService регистрирует сервис в gRPC сервере
func RegisterService(gRPCServer *grpc.Server, s Service) {
	pb.RegisterAccountsServer(gRPCServer, &s)
}
