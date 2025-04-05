package auth

import (
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	authpb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
)

func GetAuthCommands() ([]*cli.Command, error) {
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
				user := cmd.String("user")
				password := cmd.String("password")

				if user == "" || password == "" {
					fmt.Println("не указан пользователь или пароль")
					return fmt.Errorf("не указан пользователь или пароль")
				}

				request := &authpb.AuthRequest{
					Login:  user,
					Secret: password,
				}

				resp, err := authClient.Auth(ctx, request)
				if err != nil {
					fmt.Printf("ошибка аутентификации: %v", err)
					return fmt.Errorf("ошибка аутентификации: %v", err)
				}

				err = writeToken(resp.GetToken())
				if err != nil {
					return fmt.Errorf("ошибка записи токена: %v", err.Error())
				}
				fmt.Println("Авторизация прошла успешно")

				return nil
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
				user := cmd.String("user")
				password := cmd.String("password")

				if user == "" || password == "" {
					fmt.Println("не указан пользователь или пароль")
					return fmt.Errorf("не указан пользователь или пароль")
				}

				request := &authpb.RegRequest{
					Login:  user,
					Secret: password,
				}

				resp, err := authClient.Register(ctx, request)
				if err != nil {
					fmt.Printf("ошибка регистрации: %v", err)
					return err
				}

				err = writeToken(resp.GetToken())
				if err != nil {
					return fmt.Errorf("ошибка записи токена: %v", err.Error())
				}
				fmt.Println("Регистрация прошла успешно")

				return nil
			},
		},
	}, nil
}

func writeToken(token string) error {
	fileName := "token.txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	return nil
}
