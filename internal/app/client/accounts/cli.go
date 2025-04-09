package accounts

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
)

var readTokenFromFile = jwt.ReadTokenFromFile

// GetAccountCommands возвращает список команд для работы с аккаунтами
func GetAccountCommands(rootCtx context.Context) ([]*cli.Command, error) {
	cnf, err := config.MakeConfig()
	if err != nil {
		return nil, err
	}

	connAddress := cnf.ServerHost() + ":" + cnf.ServerPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	accountClient := pb.NewAccountsClient(connect)

	return []*cli.Command{
		{
			Name: "account",
			Commands: []*cli.Command{
				{
					Name:      "add",
					Usage:     "Добавление нового аккаунта",
					UsageText: "account add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return addAction(ctx, cmd, accountClient)
					},
				},
				{
					Name:      "list",
					Usage:     "Список всех аккаунтов пользователя",
					UsageText: "account list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return listAction(rootCtx, cmd, accountClient)
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление аккаунта",
					UsageText: "account rm <account_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return removeAction(ctx, cmd, accountClient)
					},
				},
				{
					Name:      "info",
					Usage:     "Детальная информация об аккаунте",
					UsageText: "account info",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return detailAction(ctx, cmd, accountClient)
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск аккаунтов",
					UsageText: "account search <url>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						return searchAction(ctx, cmd, accountClient)
					},
				},
			},
		},
	}, nil
}
