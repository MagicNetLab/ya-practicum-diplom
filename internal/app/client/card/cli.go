package card

import (
	"context"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
)

// GetCardCommands возвращает список команд для работы с картами
func GetCardCommands(rootCtx context.Context) ([]*cli.Command, error) {
	cnf, err := config.MakeConfig()
	if err != nil {
		return nil, err
	}

	connAddress := cnf.ServerHost() + ":" + cnf.ServerPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cardClient := pb.NewCardClient(connect)

	return []*cli.Command{
		{
			Name: "card",
			Commands: []*cli.Command{
				{
					Name:      "add",
					Usage:     "Добавление новой каты",
					UsageText: "card add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return addAction(rootCtx, cmd, cardClient)
					},
				},
				{
					Name:      "list",
					Usage:     "Добавление новой каты",
					UsageText: "card add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return listAction(rootCtx, cmd, cardClient)
					},
				},
				{
					Name:      "info",
					Usage:     "Информация о карте",
					UsageText: "card info <card id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return detailAction(ctx, cmd, cardClient)
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление карты",
					UsageText: "card rm <card id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return removeAction(rootCtx, cmd, cardClient)
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск карты по названию",
					UsageText: "card search <name>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return searchAction(ctx, cmd, cardClient)
					},
				},
			},
		},
	}, nil
}
