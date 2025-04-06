package note

import (
	"context"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
)

// GetNoteCommands получение команд для управления заметками
func GetNoteCommands(rootCtx context.Context) ([]*cli.Command, error) {
	cnf, err := config.MakeConfig()
	if err != nil {
		return nil, err
	}

	connAddress := cnf.ServerHost() + ":" + cnf.ServerPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	noteClient := pb.NewNoteClient(connect)

	return []*cli.Command{
		{
			Name: "note",
			Commands: []*cli.Command{
				{
					Name:      "add",
					Usage:     "Добавление заметки",
					UsageText: "note add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return addAction(rootCtx, cmd, noteClient)
					},
				},
				{
					Name:      "list",
					Usage:     "Список заметок",
					UsageText: "note list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return listAction(ctx, cmd, noteClient)
					},
				},
				{
					Name:      "info",
					Usage:     "Данные заметки",
					UsageText: "note info <noteId>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return detailAction(ctx, cmd, noteClient)
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление заметки",
					UsageText: "note rm <note_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return removeAction(ctx, cmd, noteClient)
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск заметок",
					UsageText: "note search",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return searchAction(ctx, cmd, noteClient)
					},
				},
			},
		},
	}, nil
}
