package file

import (
	"context"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

// GetFileCommands получение команд для работы с файлами.
func GetFileCommands(rootCtx context.Context) ([]*cli.Command, error) {
	cnf, err := config.MakeConfig()
	if err != nil {
		return nil, err
	}

	connAddress := cnf.ServerHost() + ":" + cnf.ServerPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	fileClient := pb.NewFilesClient(connect)

	return []*cli.Command{
		{
			Name: "file",
			Commands: []*cli.Command{
				{
					Name:      "add",
					Usage:     "Добавление файла",
					UsageText: "file add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return addAction(rootCtx, cmd, fileClient)
					},
				},
				{
					Name:      "list",
					Usage:     "Список всех файлов",
					UsageText: "file list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return listAction(rootCtx, cmd, fileClient)
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление файла",
					UsageText: "file rm <file_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return removeAction(rootCtx, cmd, fileClient)
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск по имени файла",
					UsageText: "file search",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return searchAction(rootCtx, cmd, fileClient)
					},
				},
				{
					Name:      "download",
					Usage:     "Добавление файла",
					UsageText: "file download",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return downloadAction(rootCtx, cmd, fileClient)
					},
				},
			},
		},
	}, nil
}
