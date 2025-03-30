package client

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc/metadata"
)

// accountList получает список аккаунтов пользователя и выводит его в консоль.
func accountList(ctx context.Context) {
	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := nav.Client.ListAccounts(rCtx)
	if err != nil {
		printErr("Ошибка при получении списка аккаунтов: " + err.Error())
		return
	}

	if len(res) == 0 {
		fmt.Println("\033[32m>>> Нет доступных аккаунтов.\033[0m")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Login", 24, "Url", 24, "Desc"))
	for _, r := range res {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, r.ID, 16, r.Login, 24, r.URL, 24, r.Meta))
	}
}

// accountAdd добавляет аккаунт на сайт и выводит результат в консоль.
func accountAdd(ctx context.Context) {
	var login, password, url, meta string

	fmt.Print("Введите Url сайта: ")
	_, err := fmt.Scanln(&url)
	if err != nil {
		printErr("ОШИБКА: Адрес сайта не введен.")
		return
	}

	fmt.Print("Введите Логин для сайта: ")
	_, err = fmt.Scanln(&login)
	if err != nil {
		printErr("ОШИБКА: Логин не введен.")
		return
	}

	fmt.Print("Введите Пароль для сайта: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		printErr("ОШИБКА: Логин не введен.")
		return
	}

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in := bufio.NewReader(os.Stdin)
	meta, _ = in.ReadString('\n')
	meta = strings.TrimSpace(meta)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err = nav.Client.AddAccount(rCtx, login, password, url, meta)
	if err != nil {
		printErr("Ошибка при добавлении аккаунта: " + err.Error())
		return
	}

	printInfo("Аккаунт успешно добавлен!")
}

// accountDelete удаляет аккаунт и выводит результат в консоль.
func accountDelete(ctx context.Context) {
	var id string

	fmt.Print("Введите ID аккаунта: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("ID не введен.")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	err = nav.Client.RemoveAccount(rCtx, id)
	if err != nil {
		fmt.Println("Ошибка при удалении аккаунта: " + err.Error())
		return
	}

	fmt.Println("Аккаунт успешно удален.")
}

// accountDetail получает данные аккаунта и выводит их в консоль.
func accountDetail(ctx context.Context) {
	var id string

	fmt.Print("Введите ID аккаунта: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("ID не введен.")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	account, err := nav.Client.GetAccount(rCtx, id)
	if err != nil {
		fmt.Println("Ошибка при получении аккаунта: " + err.Error())
		return
	}

	printInfo("Данные аккаунта:")
	printInfo("Url: " + account.URL)
	printInfo("Логин: " + account.Login)
	printInfo("Пароль: " + account.Password)
	printInfo("Meta: " + account.Meta)
}

// accountSearch получает список аккаунтов по поисковому запросу и выводит их в консоль.
func accountSearch(ctx context.Context) {
	var url string

	fmt.Print("Укажите адрес сайта: ")
	_, _ = fmt.Scanln(&url)

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := nav.Client.SearchAccount(rCtx, url)
	if err != nil {
		printErr("Ошибка при поиске аккаунтов: " + err.Error())
		return
	}

	if len(res) == 0 {
		fmt.Println("\033[32m>>> Ничего не найдено.\033[0m")
		return
	}

	printInfo("Результат поиска:")
	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Login", 24, "Url", 24, "Desc"))
	for _, r := range res {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, r.ID, 16, r.Login, 24, r.URL, 24, r.Meta))
	}
}
