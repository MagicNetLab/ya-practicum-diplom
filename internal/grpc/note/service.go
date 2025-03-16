package note

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"google.golang.org/grpc"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

// MakeService создание сервиса работы с заметками
func MakeService(store repository.NoteRepository, jwt config.JWTConfigurator) Service {
	return Service{store: store, jwt: jwt}
}

// Service сервис работы с заметками
type Service struct {
	pb.NoteServer
	store repository.NoteRepository
	jwt   config.JWTConfigurator
}

// Create создание заметки
func (s *Service) Create(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	note, err := models.NewNote(uid, req.GetTitle(), req.GetContent(), req.GetMeta())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.store.CreateNote(ctx, &note)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateNoteResponse{
		Note: &pb.NoteModel{
			ID:        note.GetID(),
			UID:       note.GetUID(),
			Title:     note.GetTitle(),
			Content:   note.GetContent(),
			Meta:      note.GetMeta(),
			CreatedAt: note.GetCreatedAt().Format(time.DateTime),
			UpdatedAt: note.GetUpdatedAt().Format(time.DateTime),
		},
	}, nil
}

// Get получение заметки
func (s *Service) Get(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err = uuid.Validate(req.GetID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	model, err := s.store.GetNote(ctx, req.GetID(), uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	note := &pb.NoteModel{
		ID:        model.GetID(),
		UID:       model.GetUID(),
		Title:     model.GetTitle(),
		Content:   model.GetContent(),
		Meta:      model.GetMeta(),
		CreatedAt: model.GetCreatedAt().Format(time.DateTime),
		UpdatedAt: model.GetUpdatedAt().Format(time.DateTime),
	}

	return &pb.GetNoteResponse{Note: note}, nil
}

// Update обновление заметки
func (s *Service) Update(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if uuid.Validate(req.GetID()) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	model, err := s.store.GetNote(ctx, req.GetID(), uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	err = model.SetTitle(req.GetTitle())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = model.SetContent(req.GetContent())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = model.SetMeta(req.GetMeta())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.store.UpdateNote(ctx, model)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	note := &pb.NoteModel{
		ID:        model.GetID(),
		UID:       model.GetUID(),
		Title:     model.GetTitle(),
		Content:   model.GetContent(),
		Meta:      model.GetMeta(),
		CreatedAt: model.GetCreatedAt().Format(time.DateTime),
		UpdatedAt: model.GetUpdatedAt().Format(time.DateTime),
	}

	return &pb.UpdateNoteResponse{Note: note}, nil
}

// Remove удаление заметки
func (s *Service) Remove(ctx context.Context, req *pb.RemoveNoteRequest) (*pb.RemoveNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(req.GetID()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err = s.store.RemoveNote(ctx, req.GetID(), uid)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveNoteResponse{}, nil
}

// List получение списка заметок пользователя todo deprecated
func (s *Service) List(ctx context.Context, req *pb.ListNoteRequest) (*pb.ListNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err = uuid.Validate(uid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	search := models.NoteSearch{UID: uid}
	result, err := s.store.SearchNote(ctx, &search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	notes := make([]*pb.NoteModel, 0)
	for _, row := range result {
		note := &pb.NoteModel{
			ID:        row.GetID(),
			UID:       row.GetUID(),
			Title:     row.GetTitle(),
			Content:   row.GetContent(),
			Meta:      row.GetMeta(),
			CreatedAt: row.GetCreatedAt().Format(time.DateTime),
			UpdatedAt: row.GetUpdatedAt().Format(time.DateTime),
		}
		notes = append(notes, note)
	}

	return &pb.ListNoteResponse{Notes: notes}, nil
}

// Search поиск заметок пользователя
func (s *Service) Search(ctx context.Context, req *pb.SearchNoteRequest) (*pb.SearchNoteResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := models.NoteSearch{
		UID:     uid,
		Title:   req.GetTitle(),
		Content: req.GetContent(),
		Meta:    req.GetMeta(),
		Limit:   req.GetLimit(),
		Offset:  req.GetOffset(),
	}

	result, err := s.store.SearchNote(ctx, &search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	notes := make([]*pb.NoteModel, 0)
	for _, row := range result {
		note := &pb.NoteModel{
			ID:        row.GetID(),
			UID:       row.GetUID(),
			Title:     row.GetTitle(),
			Content:   row.GetContent(),
			Meta:      row.GetMeta(),
			CreatedAt: row.GetCreatedAt().Format(time.DateTime),
			UpdatedAt: row.GetUpdatedAt().Format(time.DateTime),
		}
		notes = append(notes, note)
	}

	return &pb.SearchNoteResponse{Notes: notes}, nil
}

func RegisterService(gRPCServer *grpc.Server, s Service) {
	pb.RegisterNoteServer(gRPCServer, &s)
}
