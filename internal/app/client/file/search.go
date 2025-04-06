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

// searchFiles выполняет поиск файлов по имени и выводит информацию в консоль
func searchAction(ctx context.Context, cmd *cli.Command, fileClient pb.FilesClient) error {
	token, err := jwt.ReadTokenFromFile()
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
}
