package files

import (
	"bytes"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository/models"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/s3"
)

func MakeService(repo repository.FileRepository, storage s3.S3Client, jwt config.JWTConfigurator) Service {
	return Service{db: repo, storage: storage, jwt: jwt}
}

type Service struct {
	pb.FilesServer
	db      repository.FileRepository
	storage s3.S3Client
	jwt     config.JWTConfigurator
}

// Put добавление файла в хранилище
func (s *Service) Put(ctx context.Context, req *pb.PutFileRequest) (*pb.PutFileResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	reader := bytes.NewReader(req.GetContent())
	file, err := models.NewFile(uid, req.GetName(), req.GetMeta(), int(reader.Size()))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.db.CreateFile(ctx, file)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("error creating file: %s", err.Error()))
	}

	err = s.storage.PutObject(ctx, file.GetPath(), reader, reader.Size())
	if err != nil {
		_ = s.db.DeleteFile(ctx, file.GetID(), uid)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("error uploading file to storage: %s", err.Error()))
	}

	model := &pb.FileModel{
		Id:        file.GetID(),
		Name:      file.GetName(),
		Path:      file.GetPath(),
		Meta:      file.GetMeta(),
		Size:      int64(file.GetSize()),
		CreatedAt: time.Now().Format(time.DateTime),
	}

	return &pb.PutFileResponse{File: model}, nil
}

// List список файлов пользователя
func (s *Service) List(ctx context.Context, req *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := &models.FilesSearch{UID: uid, Limit: req.GetLimit(), Offset: req.GetOffset()}

	res, err := s.db.SearchFile(ctx, search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	files := make([]*pb.FileModel, 0)
	for _, v := range res {
		file := &pb.FileModel{
			Id:        v.GetID(),
			Name:      v.GetName(),
			Path:      v.GetPath(),
			Meta:      v.GetMeta(),
			Size:      int64(v.GetSize()),
			CreatedAt: v.GetCreatedAt().Format(time.DateTime),
		}
		files = append(files, file)
	}

	return &pb.ListFilesResponse{Files: files}, nil
}

// Download скачивание файла пользователем
func (s *Service) Download(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(req.GetId()); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id")
	}

	file, err := s.db.GetFile(ctx, req.GetId(), uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	fileContent, err := s.storage.GetObject(ctx, file.GetPath())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("error downloading file from storage: %s", err.Error()))
	}

	return &pb.DownloadFileResponse{
		Name:    file.GetName(),
		Content: fileContent,
		Size:    int64(file.GetSize()),
	}, nil
}

// Search поиск по файлам пользователя
func (s *Service) Search(ctx context.Context, req *pb.SearchFilesRequest) (*pb.SearchFilesResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	search := &models.FilesSearch{UID: uid, Limit: req.GetLimit(), Offset: req.GetOffset(), Name: req.GetName(), Meta: req.GetMeta()}
	res, err := s.db.SearchFile(ctx, search)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	files := make([]*pb.FileModel, 0)
	for _, v := range res {
		file := &pb.FileModel{
			Id:        v.GetID(),
			Name:      v.GetName(),
			Path:      v.GetPath(),
			Meta:      v.GetMeta(),
			Size:      int64(v.GetSize()),
			CreatedAt: v.GetCreatedAt().Format(time.DateTime),
		}
		files = append(files, file)
	}

	return &pb.SearchFilesResponse{Files: files}, nil
}

// Remove удаление файла пользователем
func (s *Service) Remove(ctx context.Context, req *pb.RemoveFileRequest) (*pb.RemoveFileResponse, error) {
	uid, err := jwt.GetUIDFromContext(ctx, s.jwt.GetJWTSecret())
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	if err = uuid.Validate(req.GetId()); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	file, err := s.db.GetFile(ctx, req.GetId(), uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	err = s.storage.RemoveObject(ctx, file.GetPath())
	if err != nil {
		return nil, status.Error(codes.Internal, "filed remove file")
	}

	err = s.db.DeleteFile(ctx, req.GetId(), uid)
	if err != nil {
		return nil, status.Error(codes.Internal, "filed remove file info")
	}

	return &pb.RemoveFileResponse{}, nil
}

func RegisterService(gRPCServer *grpc.Server, s Service) {
	pb.RegisterFilesServer(gRPCServer, &s)
}
