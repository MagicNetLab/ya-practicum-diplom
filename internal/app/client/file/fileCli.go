package file

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
	"strings"
)

func GetFileCommands() ([]*cli.Command, error) {
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
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						var path, name, meta string

						fmt.Print("Введите полный путь к файлу(без пробелов): ")
						_, err = fmt.Scanln(&path)
						if err != nil {
							fmt.Println("Ошибка чтения пути к файлу")
							return nil
						}

						if strings.Contains(path, " ") {
							fmt.Println("Путь к файлу не должен содержать пробелы")
							return nil
						}

						_, err = os.Stat(path)
						if err != nil {
							if os.IsNotExist(err) {
								fmt.Println("Файл не найден")
								return nil
							}

							fmt.Printf("Ошибка при получении информации о файле: " + err.Error())
							return nil
						}

						fName := strings.Split(path, "/")
						name = strings.TrimSpace(fName[len(fName)-1])

						fmt.Print("Введите название для файла без пробелов (по умолчанию: " + name + "): ")
						in := bufio.NewReader(os.Stdin)
						title, _ := in.ReadString('\n')
						if title != "\n" {
							name = strings.TrimSpace(title)
						}

						fmt.Print("Введите описание для файла: ")
						in = bufio.NewReader(os.Stdin)
						meta, _ = in.ReadString('\n')
						meta = strings.TrimSpace(meta)

						fileContent, err := readFile(path)
						if err != nil {
							return err
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.PutFileRequest{
							Name:    name,
							Meta:    meta,
							Content: fileContent,
						}

						_, err = fileClient.Put(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка при добавлении файла: " + err.Error())
							return nil
						}

						fmt.Println("Файл успешно добавлен!")

						return nil
					},
				},
				{
					Name:      "list",
					Usage:     "Список всех файлов",
					UsageText: "file list",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.ListFilesRequest{}
						res, err := fileClient.List(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка при получении списка файлов: " + err.Error())
							return nil
						}

						if len(res.Files) == 0 {
							fmt.Println("Файлы отсутствуют")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 12, "Size", 24, "Desc"))
						for _, r := range res.Files {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*d | %-*s", 38, r.GetId(), 16, r.GetName(), 12, r.GetSize(), 24, r.GetMeta()))
						}

						return nil
					},
				},
				{
					Name:      "rm",
					Usage:     "Удаление файла",
					UsageText: "file rm <file_id>",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")

							return nil
						}

						id := cmd.Args().Get(0)
						if id == "" {
							fmt.Println("Не указан id файла")

							return nil
						}

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.RemoveFileRequest{Id: id}
						_, err = fileClient.Remove(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")

								return nil
							}

							fmt.Printf("Ошибка при удалении файла: " + err.Error())

							return nil
						}

						fmt.Println("Файл успешно удален!")

						return nil
					},
				},
				{
					Name:      "search",
					Usage:     "Поиск по имени файла",
					UsageText: "file search",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						fmt.Print("Введите название файла для поиска: ")
						in := bufio.NewReader(os.Stdin)
						n, _ := in.ReadString('\n')
						name := strings.TrimSpace(n)

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.SearchFilesRequest{Name: name}
						resp, err := fileClient.Search(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка при поиске файла: " + err.Error())
							return nil
						}

						if len(resp.Files) == 0 {
							fmt.Println("Ничего не найдено")
							return nil
						}

						fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 12, "Size", 24, "Desc"))
						for _, r := range resp.Files {
							fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*d | %-*s", 38, r.GetId(), 16, r.GetName(), 12, r.GetSize(), 24, r.GetMeta()))
						}

						return nil
					},
				},
				{
					Name:      "download",
					Usage:     "Добавление файла",
					UsageText: "file download",
					Action: func(ctx context.Context, cmd *cli.Command) error {
						token, err := parseTokenFromFile()
						if err != nil {
							fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
							return nil
						}

						var id, path string
						fmt.Print("Введите ID файла: ")
						_, err = fmt.Scanln(&id)
						if err != nil {
							fmt.Println("Ошибка чтения ID файла")
							return nil
						}

						fmt.Print("Введите путь для сохранения файла (без пробелов. Пример: /path/to/file/): ")
						_, _ = fmt.Scanln(&path)

						md := metadata.Pairs("token", token)
						rCtx := metadata.NewOutgoingContext(ctx, md)
						req := &pb.DownloadFileRequest{Id: id}
						resp, err := fileClient.Download(rCtx, req)
						if err != nil {
							if status.Convert(err).Code() == codes.Unauthenticated {
								_ = os.Remove("token.txt")
								fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
								return nil
							}

							fmt.Printf("Ошибка при скачивании файла: " + err.Error())
							return nil
						}

						fPath := path + resp.GetName()
						f, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, 0600)
						if err != nil {
							fmt.Printf("Ошибка создания файла: " + err.Error())
							return nil
						}
						defer f.Close()

						writer := bufio.NewWriter(f)
						_, err = writer.Write(resp.Content)
						if err != nil {
							fmt.Printf("Ошибка записи данных в файл: " + err.Error())
							return nil
						}

						err = writer.Flush()
						if err != nil {
							fmt.Printf("Ошибка записи файла на диск: " + err.Error())
							return nil
						}

						fmt.Println("Файл успешно скачан")

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

func readFile(path string) ([]byte, error) {
	fileContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}
