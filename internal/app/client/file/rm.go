package file

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/files/proto"
)

// removeAction удаляет файл с указанным id.
func removeAction(ctx context.Context, cmd *cli.Command, fileClient pb.FilesClient) error {
	token, err := readTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")

		return nil
	}

	var id string
	fmt.Print("Введите имя карты для поиска: ")
	_, err = fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Ошибка: не заполнено имя карты")
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
}
