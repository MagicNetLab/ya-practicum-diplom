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
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// listAction получает список файлов пользователя и выводит его в консоль
func listAction(ctx context.Context, cmd *cli.Command, fileClient pb.FilesClient) error {
	token, err := jwt.ReadTokenFromFile()
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
}
