package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"google.golang.org/grpc/metadata"

	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
)

// notesList отображает список заметок пользователя
func notesList(ctx context.Context) {
	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.NotesList(rCtx)
	if err != nil {
		printErr("Ошибка при получении списка заметок: " + err.Error())
		return
	}

	if len(resp) == 0 {
		printInfo("Список заметок пуст")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s\033[0m", 38, "ID", 24, "Title", 24, "Desc"))
	for _, r := range resp {
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
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s \033[0m", 38, id, 24, title, 24, meta))
	}

}

// notesAdd добавляет новую заметку
func notesAdd(ctx context.Context) {
	data := client.NoteData{}

	fmt.Print("Введите заголовок заметки: ")
	in := bufio.NewReader(os.Stdin)
	title, _ := in.ReadString('\n')
	data.Title = strings.TrimSpace(title)

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in = bufio.NewReader(os.Stdin)
	meta, _ := in.ReadString('\n')
	data.Meta = strings.TrimSpace(meta)

	fmt.Println("Введите текст заметки (в последней строке введите --end): ")
	in = bufio.NewReader(os.Stdin)
	var input strings.Builder
	for {
		line, _ := in.ReadString('\n')
		line = strings.TrimSpace(line) // Удаление символа переноса
		if line == "--end" {
			break
		}
		input.WriteString(line + "\n")
	}
	data.Content = input.String()

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err := nav.Client.NoteCreate(rCtx, data)
	if err != nil {
		printErr("Ошибка при добавлении заметки: " + err.Error())
		return
	}

	printInfo("Заметка успешно добавлена")

}

// notesRemove удаляет заметку по ID
func notesRemove(ctx context.Context) {
	var id string
	fmt.Print("Введите ID заметки для удаления:")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("id заметки не может быть пустым")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err = nav.Client.NoteDelete(rCtx, id)
	if err != nil {
		printErr("Ошибка при удалении заметки: " + err.Error())
		return
	}

	printInfo("Заметка успешно удалена")
}

func notesDetail(ctx context.Context) {
	var id string

	fmt.Print("Введите ID заметки для просмотра: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("id заметки не может быть пустым")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	model, err := nav.Client.NoteDetail(rCtx, id)
	if err != nil {
		printErr("Ошибка при получении заметки: " + err.Error())
		return
	}

	printInfo("Заметка: " + model.ID)
	printInfo("Заголовок: " + model.Title)
	printInfo("Мета информация: " + model.Meta)
	printInfo("Содержание: \n" + model.Content)
}

// noteSearch осуществляет поиск заметок по тексту в заголовках, метах и содержании
func notesSearch(ctx context.Context) {
	var data client.NoteSearchData

	fmt.Print("Введите текст для поиска в заголовках заметок: ")
	in := bufio.NewReader(os.Stdin)
	t, _ := in.ReadString('\n')
	data.Search = strings.TrimSpace(t)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.NoteSearch(rCtx, data)
	if err != nil {
		printErr("Ошибка при поиске заметок: " + err.Error())
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s\033[0m", 38, "ID", 24, "Title", 24, "Desc"))
	for _, r := range resp {
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
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s \033[0m", 38, id, 24, title, 24, meta))
	}
}
