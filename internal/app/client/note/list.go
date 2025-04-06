package note

import (
	"context"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// listAction выводит список заметок пользователя в консоль
func listAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
	token, err := jwt.ReadTokenFromFile()
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
}
