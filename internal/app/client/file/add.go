package file

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// addAction добавляет файл на сервер
func addAction(ctx context.Context, cmd *cli.Command, fileClient pb.FilesClient) error {
	token, err := jwt.ReadTokenFromFile()
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

	fileContent, err := os.ReadFile(path)
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
}
