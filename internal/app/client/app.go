package client

import (
	"context"
	"fmt"
	appClient "github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"os/exec"
)

func Run(ctx context.Context, client appClient.AppClient) error {
	commands := initApp(ctx, client)
	err := commands.Run(ctx, os.Args)
	if err != nil {
		return err
	}

	return nil
}

var nav = navigation{Token: ""}

func initApp(ctx context.Context, client appClient.AppClient) *cli.Command {
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

						case "noteList":
							if nav.IsLoggedIn() {
								notesList(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "noteAdd":
							if nav.IsLoggedIn() {
								notesAdd(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "noteDelete":
							if nav.IsLoggedIn() {
								notesRemove(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "noteSearch":
							if nav.IsLoggedIn() {
								notesSearch(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "noteDetails":
							if nav.IsLoggedIn() {
								notesDetail(ctx)
							} else {
								printErr("Вы не авторизованы")
							}

						case "fileList":
							if nav.IsLoggedIn() {
								fileList(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "fileAdd":
							if nav.IsLoggedIn() {
								fileAdd(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "fileRemove":
							if nav.IsLoggedIn() {
								fileRemove(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "fileDownload":
							if nav.IsLoggedIn() {
								fileDownload(ctx)
							} else {
								printErr("Вы не авторизованы")
							}
						case "fileSearch":
							if nav.IsLoggedIn() {
								fileSearch(ctx)
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
			{
				Name:  "help",
				Usage: "Помощь",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					PrintHelp()
					return nil
				},
			},
		},
	}
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
