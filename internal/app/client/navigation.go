package client

import appClianet "github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"

const (
	rootCategory     = "root"
	authCategory     = "auth"
	registerCategory = "register"
	accountsCategory = "accounts"
	cardsCategory    = "cards"
	notesCategory    = "notes"
	filesCategory    = "files"
)

type navigation struct {
	Token  string
	Client appClianet.AppClient
}

func (n *navigation) IsLoggedIn() bool {
	return n.Token != ""
}

func (n *navigation) PrintAvailableCommands() {
	printInfo("Доступные команды:")

	if n.IsLoggedIn() {
		printInfo("accountList - список сохраненных аккаунтов")
		printInfo("accountAdd - добавить новый аккаунт аккаунт")
		printInfo("accountDetails- информация об аккаунте")
		printInfo("accountDelete - удалить аккаунт")
		printInfo("accountSearch - поиск аккаунта")

		printInfo("cardList - список сохраненных карт")
		printInfo("cardAdd - добавить новую карту")
		printInfo("cardDetails - информация о карте")
		printInfo("cardDelete - удалить карту")
		printInfo("cardSearch - поиск карты")

		printInfo("noteList - список сохраненных карт")
		printInfo("noteAdd - добавить новую карту")
		printInfo("noteDetails - информация о карте")
		printInfo("noteDelete - удалить карту")
		printInfo("noteSearch - поиск карты")

		printInfo("fileList - список сохраненных карт")
		printInfo("fileAdd - добавить новую карту")
		printInfo("fileDownload - скачать файл")
		printInfo("fileRemove - удалить файл")
		printInfo("fileSearch - поиск файла")

		printInfo("commands - список доступных команд")
		printInfo("clear - очистка экрана")
		printInfo("exit - выход")
	} else {
		printInfo("auth - авторизация")
		printInfo("register - регистрация")
		printInfo("commands - список доступных команд")
		printInfo("clear - очистка экрана")
		printInfo("exit - выход")
	}
}
