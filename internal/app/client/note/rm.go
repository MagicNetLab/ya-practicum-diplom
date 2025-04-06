package note

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// removeAction удаляет заметку с указанным id
func removeAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
	token, err := jwt.ReadTokenFromFile()
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
}
