package note

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/metadata"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
)

// detailAction выводит детальную информацию о заметке по ее id
func detailAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
	token, err := readTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	var id string
	fmt.Print("Введите ID заметки: ")
	_, err = fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Ошибка чтения ID заметки")
		return nil
	}

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	req := &pb.GetNoteRequest{ID: id}
	model, err := noteClient.Get(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Println("Ошибка при получении заметки: " + err.Error())
		return nil
	}
	note := model.GetNote()

	fmt.Println("Заметка: " + note.GetID())
	fmt.Println("Заголовок: " + note.GetTitle())
	fmt.Println("Мета информация: " + note.GetMeta())
	fmt.Println("Содержание: \n" + note.GetContent())

	return nil
}
