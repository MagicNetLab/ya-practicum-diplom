package app

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/interceptors"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/repository"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/s3"
)

func New(cnf config.AppConfigurator) (*Application, error) {
	repo, err := repository.NewRepository(cnf.GetDBConf())
	if err != nil {
		return nil, err
	}

	s3Client, err := s3.New(cnf.GetS3Conf())
	if err != nil {
		return nil, err
	}

	return &Application{
		repo: repo,
		cnf:  cnf,
		s3:   s3Client,
	}, nil
}

type Application struct {
	cnf    config.AppConfigurator
	server *grpc.Server
	s3     s3.S3Client
	repo   repository.Repository
}

func (app *Application) InitServer() error {
	accountService, err := account.MakeService(app.repo.GetAccountRepo(), app.cnf.GetJWTConf())
	if err != nil {
		logger.Error("failed to create account service", logger.StrArg("error", err.Error()))
		return fmt.Errorf("create account service err: %v", err)
	}

	authService, err := auth.MakeService(app.cnf.GetJWTConf(), app.repo.GetAuthRepo())
	if err != nil {
		logger.Error("failed to create auth service", logger.StrArg("error", err.Error()))
		return fmt.Errorf("creaye auth service err: %v", err)
	}

	cardService, err := card.MakeService(app.repo.GetCardRepo(), app.cnf.GetJWTConf())
	if err != nil {
		logger.Error("failed to create card service", logger.StrArg("error", err.Error()))
		return fmt.Errorf("creaye card service err: %v", err)
	}

	filesService := files.MakeService(app.repo.GetFileRepo(), app.s3, app.cnf.GetJWTConf())

	noteService := note.MakeService(app.repo.GetNoteRepo(), app.cnf.GetJWTConf())

	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.LoggerInterceptor,
			interceptors.AuthInterceptor,
			interceptors.GuestInterceptor,
		),
	}

	server := grpc.NewServer(opts...)

	account.RegisterService(server, accountService)
	auth.RegisterService(server, authService)
	card.RegisterService(server, cardService)
	files.RegisterService(server, filesService)
	note.RegisterService(server, noteService)

	app.server = server

	return nil
}

func (app *Application) Start() error {

	serverAddress := app.cnf.GetServerConf().GetHost() + ":" + app.cnf.GetServerConf().GetPort()
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		logger.Error("failed to start serv listen", logger.StrArg("error", err.Error()))
		return err
	}

	logger.Info("serv listening", logger.StrArg("address", listener.Addr().String()))

	if err := app.server.Serve(listener); err != nil {
		logger.Error("failed to serve serv", logger.StrArg("error", err.Error()))
		return err
	}

	return nil
}

func (app *Application) Stop() {
	app.server.GracefulStop()
	app.repo.Close()
}
