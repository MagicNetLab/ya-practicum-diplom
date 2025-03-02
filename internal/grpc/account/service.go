package account

import (
	"context"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"
)

// MakeService возвращает настроенный сервис аккаунтов
func MakeService(repo repository.AccountRepository) (Service, error) {
	return Service{store: repo}, nil
}

// Service  сервис аккаунтов
type Service struct {
	pb.AccountsServer
	store repository.AccountRepository
}

// GetAccount получение аккаунта по идентификатору
func (s *Service) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	id := req.GetId()
	result, err := s.store.GetAccount(ctx, id)
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

// CreateAccount создание аккаунта
func (s *Service) CreateAccount(ctx context.Context, req *pb.CreateAccountRequest) (*pb.CreateAccountResponse, error) {
	acc, err := s.store.CreateAccount(ctx, req.GetUid(), req.GetLogin(), req.GetPassword(), req.GetUrl(), req.GetDescription())
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

// RemoveAccount удаление аккаунта по идентификатору
func (s *Service) RemoveAccount(ctx context.Context, req *pb.RemoveAccountRequest) (*pb.RemoveAccountResponse, error) {
	id := req.GetId()
	err := s.store.RemoveAccount(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "account not removed")
	}

	return &pb.RemoveAccountResponse{}, nil
}

// SearchAccounts поиск аккаунтов по критериям
func (s *Service) SearchAccounts(ctx context.Context, req *pb.SearchAccountRequest) (*pb.SearchAccountResponse, error) {
	search := models.AccountSearch{}
	if req.GetUid() != "" {
		search.UID = req.GetUid()
	}

	if req.GetLogin() != "" {
		search.Login = req.GetLogin()
	}

	if req.GetUrl() != "" {
		search.URL = req.GetUrl()
	}

	if req.GetDescription() != "" {
		search.Description = req.GetDescription()
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
