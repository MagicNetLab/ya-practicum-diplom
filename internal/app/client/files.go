package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
	"github.com/jackc/pgerrcode"
	"google.golang.org/grpc/metadata"
	"os"
	"strings"
)

// Выводит список файлов пользователя
func fileList(ctx context.Context) {
	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.FileList(rCtx)
	if err != nil {
		printErr("Ошибка получения списка файлов" + err.Error())
		return
	}

	if len(resp) == 0 {
		printInfo("У вас нет сохраненных файлов")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Name", 12, "Size", 24, "Desc"))
	for _, r := range resp {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*d | %-*s\033[0m", 38, r.ID, 16, r.Name, 12, r.Size, 24, r.Meta))
	}
}

// Добавляет файл в хранилище пользователя
func fileAdd(ctx context.Context) {
	var data = client.FileData{}

	fmt.Print("Введите полный путь к файлу: ")
	_, err := fmt.Scanln(&data.Path)
	if err != nil {
		printErr("Ошибка чтения пути к файлу (без пробелов)")
		return
	}

	if strings.Contains(data.Path, " ") {
		printErr("Путь к файлу не должен содержать пробелы")
		return
	}

	_, err = os.Stat(data.Path)
	if err != nil {
		if os.IsNotExist(err) {
			printErr("Файл не найден")
			return
		}

		printErr("Ошибка чтения пути к файлу: " + err.Error())
		return
	}

	fName := strings.Split(data.Path, "/")
	data.Name = strings.TrimSpace(fName[len(fName)-1])

	fmt.Print("Введите название для файла без пробелов (по умолчанию: " + fName[len(fName)-1] + "): ")
	in := bufio.NewReader(os.Stdin)
	title, _ := in.ReadString('\n')
	if title != "\n" {
		data.Name = strings.TrimSpace(title)
	}

	fmt.Print("Введите описание для файла: ")
	in = bufio.NewReader(os.Stdin)
	meta, _ := in.ReadString('\n')
	data.Meta = strings.TrimSpace(meta)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err = nav.Client.FileAdd(rCtx, data)
	if err != nil {
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			printErr("Файл с таким именем уже существует")
			return
		}
		printErr("Ошибка добавления файла" + err.Error())
		return
	}

	printInfo("Файл успешно добавлен")
}

// Удаляет файл из хранилища пользователя
func fileRemove(ctx context.Context) {
	var id string
	fmt.Print("Введите ID файла: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("ID файла не может быть пустым")
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err = nav.Client.FileRemove(rCtx, id)
	if err != nil {
		printErr("Ошибка удаления файла" + err.Error())
		return
	}
	printInfo("Файл успешно удален")
}

func fileDownload(ctx context.Context) {
	var id, path string
	fmt.Print("Введите ID файла: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("ID файла не может быть пустым")
		return
	}

	fmt.Print("Введите путь для сохранения файла (без пробелов. Пример: /path/to/file/): ")
	_, _ = fmt.Scanln(&path)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	file, err := nav.Client.FileDownload(rCtx, id)
	if err != nil {
		printErr("Ошибка скачивания файла" + err.Error())
		return
	}

	fPath := path + file.Name
	f, err := os.OpenFile(fPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		printErr("Ошибка создания файла: " + err.Error())
		return
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	_, err = writer.Write(file.Content)
	if err != nil {
		printErr("Ошибка записи данных в файл: " + err.Error())
		return
	}

	err = writer.Flush()
	if err != nil {
		printErr("Ошибка записи файла на диск: " + err.Error())
		return
	}

	printInfo("Файл успешно скачан")
}

func fileSearch(ctx context.Context) {
	var name, meta string

	fmt.Print("Введите название файла для поиска: ")
	in := bufio.NewReader(os.Stdin)
	n, _ := in.ReadString('\n')
	name = strings.TrimSpace(n)

	fmt.Print("Введите описание файла для поиска: ")
	in = bufio.NewReader(os.Stdin)
	m, _ := in.ReadString('\n')
	meta = strings.TrimSpace(m)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.FileSearch(rCtx, name, meta)
	if err != nil {
		printErr("Ошибка получения списка файлов" + err.Error())
		return
	}

	if len(resp) == 0 {
		printInfo("Нет файлов с такими параметрами")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Name", 12, "Size", 24, "Desc"))
	for _, r := range resp {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*d | %-*s\033[0m", 38, r.ID, 16, r.Name, 12, r.Size, 24, r.Meta))
	}

}
