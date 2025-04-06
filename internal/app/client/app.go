package client

import (
	"context"
	"os"
	"time"

	"github.com/urfave/cli/v3"

	accountCli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client/accounts"
	authCli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client/auth"
	cardCli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client/card"
	fileCli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client/file"
	noteCli "github.com/MagicNetLab/ya-practicum-diplom/internal/app/client/note"
)

// Run Запуск клиента приложения
func Run(ctx context.Context) error {
	commands, err := initApp(ctx)
	if err != nil {
		return err
	}

	err = commands.Run(ctx, os.Args)
	if err != nil {
		return err
	}

	return nil
}

// initApp Инициализация  клиента приложения
func initApp(ctx context.Context) (*cli.Command, error) {
	clientVersion := os.Getenv("CLIENT_VERSION")
	buildDate := time.Now().Format("02.01.2006 15:04:05")

	commands, err := buildCommands(ctx)
	if err != nil {
		return nil, err
	}

	return &cli.Command{
		Name:     "GophKeeper",
		Usage:    "GophKeeper client",
		Version:  clientVersion + " / build: " + buildDate,
		Commands: commands,
	}, nil
}

// buildCommands Создание команд для приложения
func buildCommands(ctx context.Context) ([]*cli.Command, error) {
	commands := make([]*cli.Command, 0)

	// Авторизация/регистрация
	authCommands, err := authCli.GetAuthCommands(ctx)
	if err != nil {
		return nil, err
	}
	commands = append(commands, authCommands...)

	// Работа с аккаунтами
	accountCommands, err := accountCli.GetAccountCommands(ctx)
	if err != nil {
		return nil, err
	}
	commands = append(commands, accountCommands...)

	// Работа с картами
	cardCommands, err := cardCli.GetCardCommands(ctx)
	if err != nil {
		return nil, err
	}
	commands = append(commands, cardCommands...)

	// Работа с заметками
	notesCommands, err := noteCli.GetNoteCommands(ctx)
	if err != nil {
		return nil, err
	}
	commands = append(commands, notesCommands...)

	fileCommands, err := fileCli.GetFileCommands(ctx)
	if err != nil {
		return nil, err
	}
	commands = append(commands, fileCommands...)

	return commands, nil
}
