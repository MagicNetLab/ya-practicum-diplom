package client

import (
	"context"
	"fmt"
	appClianet "github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"os/exec"
)

func Run(ctx context.Context, client appClianet.AppClient) error {
	commands := initApp(ctx, client)
	err := commands.Run(ctx, os.Args)
	if err != nil {
		return err
	}

	return nil
}

var nav = navigation{Token: ""}
var done = make(chan struct{})

func initApp(ctx context.Context, client appClianet.AppClient) *cli.Command {
	nav.Client = client

	return &cli.Command{
		Name:    "GophKeeper",
		Usage:   "GophKeeper client",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Запуск приложение",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					printInfo("Приложение запущено")
					nav.PrintAvailableCommands()

					for {
						var input string
						fmt.Print("Введите команду: ")
						_, _ = fmt.Scanln(&input)
						switch input {
						case "auth":
							if !nav.IsLoggedIn() {
								printInfo("Авторизация")
								_ = auth(ctx)
							}
						case "register":
							if !nav.IsLoggedIn() {
								printInfo("Регистрация")
								_ = register(ctx)
							}
						case "exit":
							printInfo("Приложение закрыто")
							return nil
						case "clear":
							cmd := exec.Command("clear")
							cmd.Stdout = os.Stdout
							err := cmd.Run()
							if err != nil {
								log.Println(err.Error())
							}
						case "commands":
							nav.PrintAvailableCommands()
						case "accountList":
							if nav.IsLoggedIn() {
								accountList(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "accountAdd":
							if nav.IsLoggedIn() {
								accountAdd(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "accountDelete":
							if nav.IsLoggedIn() {
								accountDelete(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "accountSearch":
							if nav.IsLoggedIn() {
								accountSearch(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "accountDetails":
							if nav.IsLoggedIn() {
								accountDetail(ctx)
							} else {
								printErr("Вы не авторизованы")
							}

						case "cardList":
							if nav.IsLoggedIn() {
								cardsList(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "cardAdd":
							if nav.IsLoggedIn() {
								cardsCreate(ctx)
							} else {
								cardsDelete(ctx)
							}
						case "cardDelete":
							if nav.IsLoggedIn() {
								cardsDelete(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "cardSearch":
							if nav.IsLoggedIn() {
								cardsSearch(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "cardDetails":
							if nav.IsLoggedIn() {
								cardDetail(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "commandList":
							if nav.IsLoggedIn() {
								cardsSearch(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						default:
							printErr("Команда не найдена")
						}
					}
				},
			},
		},
	}
}

// auth выполняет авторизацию пользователя
func auth(ctx context.Context) error {
	var login string
	var password string

	fmt.Print("Введите логин: ")
	_, err := fmt.Scanln(&login)
	if err != nil {
		printErr("ОШИБКА: Логин не введен")
		return err
	}
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}

	jwtToken, err := nav.Client.Auth(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			printErr("Ошибка аутентификации: " + e.Message())
			return err
		}
		printErr("Ошибка аутентификации")
		return err
	}

	printSuccess("Аутентификация прошла успешно")
	nav.Token = jwtToken
	clearScreen()
	nav.PrintAvailableCommands()

	return nil
}

func register(ctx context.Context) error {
	var login string
	var password string
	var rePassword string

	fmt.Print("Введите логин: ")
	_, err := fmt.Scanln(&login)
	if err != nil {
		printErr("Ошибка: логин не введен")
		return err
	}
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}
	fmt.Print("Введите пароль еще раз: ")
	_, err = fmt.Scanln(&rePassword)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}

	if password != rePassword {
		printErr("ОШИБКА: Пароли не совпадают")
		return fmt.Errorf("пароли не совпадают")
	}

	jwtToken, err := nav.Client.Register(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			printErr("Ошибка регистрации: " + e.Message())
			return err
		}
		printErr("Ошибка регистрации")
		return err
	}

	clearScreen()

	printSuccess("Пользователь успешно зарегистрирован")
	nav.PrintAvailableCommands()
	nav.Token = jwtToken
	return nil
}

func printErr(str string) {
	fmt.Printf("\033[31m>>> %s\033[0m\n", str)
}

func printInfo(str string) {
	fmt.Printf("\033[32m>>> %s\033[0m\n", str)
}

func printSuccess(str string) {
	fmt.Printf("\033[32m>>> %s\033[0m\n", str)
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
