package accounts

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

func GetAccountCommands() ([]*cli.Command, error) {
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
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return err
						}

						var login, password, url, meta string

						fmt.Print("Введите Url сайта: ")
						_, err = fmt.Scanln(&url)
						if err != nil {
							fmt.Println("ОШИБКА: Адрес сайта не введен.")
							return err
						}

						fmt.Print("Введите Логин для сайта: ")
						_, err = fmt.Scanln(&login)
						if err != nil {
							fmt.Println("ОШИБКА: Логин не введен.")
							return nil
						}

						fmt.Print("Введите Пароль для сайта: ")
						_, err = fmt.Scanln(&password)
						if err != nil {
							fmt.Println("ОШИБКА: Логин не введен.")
							return nil
						}

						fmt.Print("Введите дополнительную информацию (не обязательно): ")
						in := bufio.NewReader(os.Stdin)
						meta, _ = in.ReadString('\n')
						meta = strings.TrimSpace(meta)

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						req := &pb.CreateAccountRequest{
							Login:       login,
							Password:    password,
							Url:         url,
							Description: meta,
						}

						_, err = accountClient.Create(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
							}
							return nil
						}

						fmt.Println("Аккаунт успешно добавлен!")

						return nil
					},
				},
				{
					Name:      "list",
					Usage:     "Список всех аккаунтов пользователя",
					UsageText: "account list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return err
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						req := &pb.SearchAccountRequest{}
						resp, err := accountClient.Search(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Println("Ошибка при получении списка аккаунтов.")
							return nil
						}

						if len(resp.GetAcc()) == 0 {
							fmt.Println("Список аккаунтов пуст.")
							return nil
						}

						fmt.Println("--- Список аккаунтов пользователя ---")
						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Login", 24, "Url", 24, "Desc"))
						for _, row := range resp.GetAcc() {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, row.GetId(), 16, row.GetLogin(), 24, row.GetUrl(), 24, row.GetDescription()))
						}

						return nil
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление аккаунта",
					UsageText: "account rm <account_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Необходимо указать id аккаунта")
							return nil
						}
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались.")
							return err
						}
						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						req := &pb.RemoveAccountRequest{Id: id}

						_, err = accountClient.Remove(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Println("Ошибка при удалении аккаунта.")
							return nil
						}

						fmt.Println("Аккаунт успешно удален!")
						return nil
					},
				},
				{
					Name:      "info",
					Usage:     "Детальная информация об аккаунте",
					UsageText: "account info <account_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Необходимо указать id аккаунта")
							return nil
						}

						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return err
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.GetAccountRequest{Id: id}

						resp, err := accountClient.Get(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}
							fmt.Println("Ошибка при получении информации об аккаунте.")
							return nil
						}

						acc := resp.GetAccount()
						fmt.Println("Данные аккаунта:")
						fmt.Println("Url: " + acc.GetUrl())
						fmt.Println("Логин: " + acc.GetLogin())
						fmt.Println("Пароль: " + acc.GetPassword())
						fmt.Println("Meta: " + acc.GetDescription())

						return nil
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск аккаунтов",
					UsageText: "account search <url>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						search := cmd.Args().Get(0)
						if search == "" {
							fmt.Println("Необходимо указать адрес сайта.")
							return nil
						}
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return err
						}
						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						req := &pb.SearchAccountRequest{Search: search}
						resp, err := accountClient.Search(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}
							fmt.Println("Ошибка при получении списка аккаунтов.")
							return nil
						}

						if len(resp.GetAcc()) == 0 {
							fmt.Println("Ни чего не найдено.")
							return nil
						}

						fmt.Println("--- Результат поиска ---")
						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Login", 24, "Url", 24, "Desc"))
						for _, row := range resp.GetAcc() {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, row.GetId(), 16, row.GetLogin(), 24, row.GetUrl(), 24, row.GetDescription()))
						}

						return nil
					},
				},
			},
		},
	}, nil
}

func parseTokenFromFile() (string, error) {
	fileName := "token.txt"
	token := ""
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		token = strings.TrimSpace(scanner.Text())
		break
	}

	if token == "" {
		return "", fmt.Errorf("token not found")
	}

	return token, nil
}
