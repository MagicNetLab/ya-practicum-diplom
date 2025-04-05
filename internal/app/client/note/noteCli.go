package note

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strings"
	"unicode/utf8"
)

func GetNoteCommands() ([]*cli.Command, error) {
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
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						var title, meta, content string

						fmt.Print("Введите заголовок заметки: ")
						in := bufio.NewReader(os.Stdin)
						title, _ = in.ReadString('\n')
						title = strings.TrimSpace(title)

						fmt.Print("Введите дополнительную информацию (не обязательно): ")
						in = bufio.NewReader(os.Stdin)
						meta, _ = in.ReadString('\n')
						meta = strings.TrimSpace(meta)

						fmt.Println("Введите текст заметки (в последней строке введите --end): ")
						in = bufio.NewReader(os.Stdin)
						var input strings.Builder
						for {
							line, _ := in.ReadString('\n')
							line = strings.TrimSpace(line)
							if line == "--end" {
								break
							}
							input.WriteString(line + "\n")
						}
						content = input.String()

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						_, err = noteClient.Create(rCtx, &pb.CreateNoteRequest{
							Title:   title,
							Meta:    meta,
							Content: content,
						})

						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}
							fmt.Printf("Ошибка при создании заметки: %s", err)
							return nil
						}

						fmt.Println("Заметка успешно добавлена")

						return nil
					},
				},
				{
					Name:      "list",
					Usage:     "Список заметок",
					UsageText: "note list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)

						resp, err := noteClient.List(rCtx, &pb.ListNoteRequest{})
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка при получении списка заметок: %s", err)
							return nil
						}

						if len(resp.GetNotes()) == 0 {
							fmt.Println("Ничего не найдено")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s", 38, "ID", 24, "Title", 24, "Desc"))
						for _, r := range resp.GetNotes() {
							id := r.GetID()
							title := r.GetTitle()
							meta := r.GetMeta()
							if utf8.RuneCountInString(id) > 38 {
								id = string([]rune(id)[:38]) + "..."
							}
							if utf8.RuneCountInString(title) > 24 {
								title = string([]rune(title)[:21]) + "..."
							}
							if utf8.RuneCountInString(meta) > 24 {
								meta = string([]rune(meta)[:21]) + "..."
							}
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s", 38, id, 24, title, 24, meta))
						}

						return nil
					},
				},
				{
					Name:      "info",
					Usage:     "Данные заметки",
					UsageText: "note info <noteId>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Необходимо указать id заметки")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.GetNoteRequest{ID: id}
						model, err := noteClient.Get(rCtx, req)
						if err != nil {
							fmt.Println("Ошибка при получении заметки: " + err.Error())
							return nil
						}
						note := model.GetNote()

						fmt.Println("Заметка: " + note.GetID())
						fmt.Println("Заголовок: " + note.GetTitle())
						fmt.Println("Мета информация: " + note.GetMeta())
						fmt.Println("Содержание: \n" + note.GetContent())

						return nil
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление заметки",
					UsageText: "note rm <note_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Необходимо указать id заметки")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.RemoveNoteRequest{ID: id}
						_, err = noteClient.Remove(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Println("Ошибка при удалении заметки: " + err.Error())
							return nil
						}

						fmt.Println("Заметка удалена")

						return nil
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск заметок",
					UsageText: "note search",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						fmt.Print("Введите текст для поиска в заголовках заметок: ")
						in := bufio.NewReader(os.Stdin)
						t, _ := in.ReadString('\n')
						search := strings.TrimSpace(t)

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.SearchNoteRequest{Search: search}
						resp, err := noteClient.Search(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Println("Ошибка при получении заметок: " + err.Error())
							return nil
						}

						if len(resp.GetNotes()) == 0 {
							fmt.Println("Ничего не найдено")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s", 38, "ID", 24, "Title", 24, "Desc"))
						for _, r := range resp.GetNotes() {
							id := r.ID
							title := r.Title
							meta := r.Meta
							if utf8.RuneCountInString(id) > 38 {
								id = string([]rune(id)[:38]) + "..."
							}
							if utf8.RuneCountInString(title) > 24 {
								title = string([]rune(title)[:21]) + "..."
							}
							if utf8.RuneCountInString(meta) > 24 {
								meta = string([]rune(meta)[:21]) + "..."
							}
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s", 38, id, 24, title, 24, meta))
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
