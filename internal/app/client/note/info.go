package note

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/metadata"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// detailAction выводит детальную информацию о заметке по ее id
func detailAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
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
}
