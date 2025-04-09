package note

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

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/note/proto"
)

// addAction добавляет заметку в хранилище
func addAction(ctx context.Context, cmd *cli.Command, noteClient pb.NoteClient) error {
	token, err := readTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	var title, meta, content string

	fmt.Print("Введите заголовок заметки: ")
	in := bufio.NewReader(os.Stdin)
	title, _ = in.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in = bufio.NewReader(os.Stdin)
	meta, _ = in.ReadString('\n')
	meta = strings.TrimSpace(meta)

	fmt.Println("Введите текст заметки (в последней строке введите --end): ")
	in = bufio.NewReader(os.Stdin)
	var input strings.Builder
	for {
		line, _ := in.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "--end" {
			break
		}
		input.WriteString(line + "\n")
	}
	content = input.String()

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	_, err = noteClient.Create(rCtx, &pb.CreateNoteRequest{
		Title:   title,
		Meta:    meta,
		Content: content,
	})

	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Printf("Ошибка при создании заметки: %s", err)
		return nil
	}

	fmt.Println("Заметка успешно добавлена")

	return nil
}
