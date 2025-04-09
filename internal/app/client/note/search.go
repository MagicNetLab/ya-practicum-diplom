package note

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
)

// searchAction осуществляет поиск заметок по заголовку.
func searchAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
	token, err := readTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	fmt.Print("Введите текст для поиска в заголовках заметок: ")
	in := bufio.NewReader(os.Stdin)
	t, _ := in.ReadString('\n')
	search := strings.TrimSpace(t)

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	req := &pb.SearchNoteRequest{Search: search}
	resp, err := noteClient.Search(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}

		fmt.Println("Ошибка при получении заметок: " + err.Error())
		return nil
	}

	if len(resp.GetNotes()) == 0 {
		fmt.Println("Ничего не найдено")
		return nil
	}

	fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s", 38, "ID", 24, "Title", 24, "Desc"))
	for _, r := range resp.GetNotes() {
		id := r.ID
		title := r.Title
		meta := r.Meta
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
