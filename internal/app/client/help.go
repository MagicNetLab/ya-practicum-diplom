package client

import "fmt"

// PrintHelp - выводит список доступных команд
func PrintHelp() {
	fmt.Println("Доступные команды:")
	fmt.Println("  auth                  - авторизация пользователя")
	fmt.Println("  register              - регистрация нового пользователя")
	fmt.Println("  ")
	fmt.Println("  commands              - список доступных команд")
	fmt.Println("  clear                 - очистка экрана")
	fmt.Println("  exit                  - выход из приложения")
	fmt.Println("  ")
	fmt.Println("  accountList           - список сохраненных аккаунтов")
	fmt.Println("  accountAdd            - добавить новый аккаунт аккаунт")
	fmt.Println("  accountDetails        - информация об аккаунте")
	fmt.Println("  accountDelete         - удалить аккаунт")
	fmt.Println("  accountSearch         - поиск аккаунта")
	fmt.Println("  ")
	fmt.Println("  cardList              - список сохраненных карт")
	fmt.Println("  cardAdd               - добавить новую карту")
	fmt.Println("  cardDetails           - информация о карте")
	fmt.Println("  cardDelete            - удалить карту")
	fmt.Println("  cardSearch            - поиск карты")
	fmt.Println("  ")
	fmt.Println("  noteList              - список сохраненных карт")
	fmt.Println("  noteAdd               - добавить новую карту")
	fmt.Println("  noteDetails           - информация о карте")
	fmt.Println("  noteDelete            - удалить карту")
	fmt.Println("  noteSearch            - поиск карты")
	fmt.Println("  ")
	fmt.Println("  fileList              - список сохраненных карт")
	fmt.Println("  fileAdd               - добавить новую карту")
	fmt.Println("  fileDownload          - скачать файл")
	fmt.Println("  fileRemove            - удалить файл")
	fmt.Println("  fileSearch            - поиск файла")
}
