package note

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
)

func MakeService(store repository.NoteRepository) *Service {
	return &Service{store: store}
}

// Service сервис работы с заметками
type Service struct {
	pb    pb.NoteServer
	store repository.NoteRepository
}

// CreateNote создание заметки
func (s *Service) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	note, err := models.NewNote(req.GetUID(), req.GetTitle(), req.GetContent(), req.GetMeta())
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

// GetNote получение заметки
func (s *Service) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	err := uuid.Validate(req.GetID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	model, err := s.store.GetNote(ctx, req.GetID())
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

// UpdateNote обновление заметки
func (s *Service) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	if uuid.Validate(req.GetID()) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}

	if uuid.Validate(req.GetUID()) != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uid")
	}

	model, err := s.store.GetNote(ctx, req.GetUID())
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

// RemoveNote удаление заметки
func (s *Service) RemoveNote(ctx context.Context, req *pb.RemoveNoteRequest) (*pb.RemoveNoteResponse, error) {
	if err := uuid.Validate(req.GetID()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.store.RemoveNote(ctx, req.GetID())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveNoteResponse{}, nil
}

// ListNotes получение списка заметок пользователя
func (s *Service) ListNotes(ctx context.Context, req *pb.ListNoteRequest) (*pb.ListNoteResponse, error) {
	err := uuid.Validate(req.GetUID())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	search := models.NoteSearch{UID: req.GetUID()}
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

// SearchNotes поиск заметок пользователя
func (s *Service) SearchNotes(ctx context.Context, req *pb.SearchNoteRequest) (*pb.SearchNoteResponse, error) {
	if err := uuid.Validate(req.GetUID()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	search := models.NoteSearch{
		UID:     req.GetUID(),
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
