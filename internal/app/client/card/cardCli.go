package card

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strconv"
	"strings"
)

func GetCardCommands() ([]*cli.Command, error) {
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
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						var name, number, month, year, cvc, pin, meta string
						fmt.Print("Введите название карты: ")
						_, err = fmt.Scan(&name)
						if err != nil {
							fmt.Println("Ошибка: не заполнено название карты")
							return nil
						}

						fmt.Print("Введите номер карты: ")
						_, err = fmt.Scan(&number)
						if err != nil {
							fmt.Println("Ошибка: не заполнен номер карты")
							return nil
						}

						fmt.Print("Введите месяц окончания действия карты: ")
						_, err = fmt.Scan(&month)
						if err != nil {
							fmt.Println("Ошибка: не заполнен месяц")
							return nil
						}
						cardMonth, err := strconv.Atoi(month)
						if err != nil {
							fmt.Println("Ошибка: не заполнен месяц")
							return nil
						}

						fmt.Print("Введите год окончания действия карты: ")
						_, err = fmt.Scan(&year)
						if err != nil {
							fmt.Println("Ошибка: не заполнен год")
							return nil
						}
						cardYear, err := strconv.Atoi(year)
						if err != nil {
							fmt.Println("Ошибка: не заполнен год")
							return nil
						}

						fmt.Print("Введите CVC карты: ")
						_, err = fmt.Scan(&cvc)
						if err != nil {
							fmt.Println("Ошибка: не заполнен CVC")
							return nil
						}

						fmt.Print("Введите PIN карты: ")
						_, err = fmt.Scan(&pin)
						if err != nil {
							fmt.Println("Ошибка: не заполнен PIN")
							return nil
						}

						fmt.Print("Введите дополнительную информацию (не обязательно): ")
						in := bufio.NewReader(os.Stdin)
						meta, _ = in.ReadString('\n')
						if meta == "\n" {
							meta = ""
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						_, err = cardClient.Create(rCtx, &pb.CreateCardRequest{
							Name:   name,
							Number: number,
							Meta:   meta,
							Month:  int32(cardMonth),
							Year:   int32(cardYear),
							CVC:    cvc,
							PIN:    pin,
						})
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							if status.Convert(err).Code() == codes.InvalidArgument {
								fmt.Println("Карта не добавлена. Проверьте правильность заполнения данных.")
								return nil
							}

							fmt.Printf("Ошибка при добавлении карты: %v\n", err)
							return nil
						}

						fmt.Println("Карта успешно добавлена")
						return nil
					},
				},
				{
					Name:      "list",
					Usage:     "Добавление новой каты",
					UsageText: "card add",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}
						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						req := &pb.ListCardRequest{}
						resp, err := cardClient.List(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка получения списка карт: %v\n", err)
							return nil
						}

						if len(resp.GetCards()) == 0 {
							fmt.Println("У вас нет карт")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
						for _, row := range resp.GetCards() {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, row.GetID(), 16, row.GetName(), 24, row.GetNumber(), 24, row.GetMeta()))
						}

						return nil
					},
				},
				{
					Name:      "info",
					Usage:     "Информация о карте",
					UsageText: "card info <card id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Ошибка: не указан ID карты")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						res, err := cardClient.Get(rCtx, &pb.GetCardRequest{ID: id})
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Println("Ошибка получения подробной информации о карте: " + err.Error())
							return nil
						}

						card := res.GetCard()

						fmt.Println("Информация о карте:")
						fmt.Println(fmt.Sprintf("Наименование: %s", card.GetName()))
						fmt.Println(fmt.Sprintf("Номер: %s", card.GetNumber()))
						fmt.Println(fmt.Sprintf("Срок действия: %d/%d", card.GetMonth(), card.GetYear()))
						fmt.Println(fmt.Sprintf("PIN: %s", card.GetPIN()))
						fmt.Println(fmt.Sprintf("CVC: %s", card.GetCVC()))
						fmt.Println(fmt.Sprintf("Дополнительная информация: %s", card.GetMeta()))

						return nil
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление карты",
					UsageText: "card rm <card id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Ошибка: не указан ID карты")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						_, err = cardClient.Delete(rCtx, &pb.DeleteCardRequest{ID: id})
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}
							fmt.Println("Ошибка удаления карты: " + err.Error())
							return nil
						}

						fmt.Println("Карта успешно удалена")

						return nil
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск карты по названию",
					UsageText: "card search <name>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						var search string
						fmt.Print("Введите имя карты для поиска: ")
						_, err = fmt.Scanln(&search)
						if err != nil {
							fmt.Println("Ошибка: не заполнено имя карты")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.SearchCardRequest{Name: search}
						res, err := cardClient.Search(rCtx, req)

						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}
							fmt.Println("Ошибка получения списка карт: " + err.Error())
							return nil
						}

						if len(res.GetCards()) == 0 {
							fmt.Println("Ничего не найдено")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
						for _, r := range res.GetCards() {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, r.GetID(), 16, r.GetName(), 24, r.GetNumber(), 24, r.GetMeta()))
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
