package auth

import (
	"context"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	authpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetAuthCommands - возвращает список команд аутентификации/регистрации пользователя
func GetAuthCommands(rootCtx context.Context) ([]*cli.Command, error) {
	cnf, err := config.MakeConfig()
	if err != nil {
		return nil, err
	}

	connAddress := cnf.ServerHost() + ":" + cnf.ServerPort()
	connect, err := grpc.NewClient(connAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	authClient := authpb.NewAuthClient(connect)

	return []*cli.Command{
		{
			Name:      "auth",
			Usage:     "Аутентификация пользователя",
			UsageText: "auth -u <user> -p <password>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "user",
					Aliases:  []string{"u"},
					Required: true,
					Usage:    "Имя пользователя",
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "Пароль",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return authAction(rootCtx, cmd, authClient)
			},
		},
		{
			Name:      "register",
			Usage:     "Регистрация нового пользователя",
			UsageText: "register -u <user> -p <password>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "user",
					Aliases:  []string{"u"},
					Required: true,
					Usage:    "Имя пользователя",
				},
				&cli.StringFlag{
					Name:     "password",
					Aliases:  []string{"p"},
					Required: true,
					Usage:    "Пароль",
				},
			},
			Action: func(ctx context.Context, cmd *cli.Command) error {
				return registerAction(rootCtx, cmd, authClient)
			},
		},
	}, nil
}
