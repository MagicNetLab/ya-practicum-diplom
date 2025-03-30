package card

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/google/uuid"
)

// MakeService возвращает настроенный сервис карт
func MakeService(repo repository.CardRepository, jwtConfig config.JWTConfigurator) (Service, error) {
	return Service{store: repo, jwt: jwtConfig}, nil
}

// Service  сервис карт
type Service struct {
	pb.CardServer
	store repository.CardRepository
	jwt   config.JWTConfigurator
}

// Get получение карты по идентификатору
func (s *Service) Get(ctx context.Context, req *pb.GetCardRequest) (*pb.GetCardResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(uid); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err := uuid.Validate(req.ID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	card, err := s.store.GetCardByID(ctx, req.ID, uid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "card not found")
	}

	model := pb.CardModel{
		ID:     card.GetID(),
		Name:   card.GetName(),
		Number: card.GetNumber(),
		Month:  int32(card.GetMonth()),
		Year:   int32(card.GetYear()),
		Meta:   card.GetMeta(),
		CVC:    card.GetCVC(),
		PIN:    card.GetPIN(),
	}

	return &pb.GetCardResponse{Card: &model}, nil
}

// Create создание карты
func (s *Service) Create(ctx context.Context, req *pb.CreateCardRequest) (*pb.CreateCardResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(uid); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	card := models.Card{
		ID:        uuid.New().String(),
		UID:       uid,
		Name:      req.Name,
		Meta:      req.Meta,
		Number:    req.Number,
		Month:     int(req.Month),
		Year:      int(req.Year),
		CVC:       req.CVC,
		PIN:       req.PIN,
		CreatedAt: time.Now(),
	}
	card.MaskNumber()

	err = card.Validate()
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid card data")
	}

	err = s.store.CreateCard(ctx, &card)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &pb.CreateCardResponse{Card: &pb.ShortCardModel{
		ID:     card.GetID(),
		Name:   card.GetName(),
		Number: card.GetMask(),
		Meta:   card.GetMeta(),
	}}, nil
}

// Delete удаление карты по идентификатору
func (s *Service) Delete(ctx context.Context, req *pb.DeleteCardRequest) (*pb.DeleteCardResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(uid); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err := uuid.Validate(req.ID); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	res := s.store.DeleteCard(ctx, req.ID, uid)
	if res != nil {
		return nil, status.Errorf(codes.NotFound, "card not found")
	}

	return &pb.DeleteCardResponse{}, nil
}

// Search поиск карт по запросу
func (s *Service) Search(ctx context.Context, req *pb.SearchCardRequest) (*pb.SearchCardResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(uid); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := models.CardSearch{UID: uid, Name: req.Name, Limit: int(req.Limit), Offset: int(req.Offset)}
	res, err := s.store.SearchCards(ctx, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search cards")
	}

	cards := make([]*pb.ShortCardModel, 0)
	for _, card := range res {
		model := pb.ShortCardModel{
			ID:     card.GetID(),
			Name:   card.GetName(),
			Number: card.GetMask(),
			Meta:   card.GetMeta(),
		}

		cards = append(cards, &model)
	}

	return &pb.SearchCardResponse{Cards: cards}, nil
}

// List получение списка карт по идентификатору владельца
func (s *Service) List(ctx context.Context, req *pb.ListCardRequest) (*pb.ListCardResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(uid); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := models.CardSearch{UID: uid}
	res, err := s.store.SearchCards(ctx, search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "faled to get cards list")
	}

	cards := make([]*pb.ShortCardModel, 0, len(res))
	for _, card := range res {
		model := pb.ShortCardModel{
			ID:     card.GetID(),
			Name:   card.GetName(),
			Number: card.GetMask(),
			Meta:   card.GetMeta(),
		}

		cards = append(cards, &model)
	}

	return &pb.ListCardResponse{Cards: cards}, nil
}

// RegisterService регистрирует сервис карт
func RegisterService(gRPCServer *grpc.Server, s Service) {
	pb.RegisterCardServer(gRPCServer, &s)
}
