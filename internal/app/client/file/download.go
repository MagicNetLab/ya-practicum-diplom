package file

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

// downloadAction скачивает файл по его ID на указанный путь на диске
func downloadAction(ctx context.Context, cmd *cli.Command, fileClient pb.FilesClient) error {
	token, err := readTokenFromFile()
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
}
