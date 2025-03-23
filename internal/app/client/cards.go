package client

import (
	"bufio"
	"context"
	"fmt"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/services/client"
	"google.golang.org/grpc/metadata"
	"os"
)

// cardsList выводит список карт клиента на экран
func cardsList(ctx context.Context) {
	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.ListCards(rCtx)
	if err != nil {
		printErr("Ошибка получения списка карт")
		return
	}

	if len(resp) == 0 {
		printInfo("У вас нет карт")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
	for _, r := range resp {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, r.ID, 16, r.Name, 24, r.Number, 24, r.Meta))
	}
}

// cardsCreate создаёт карту клиента
func cardsCreate(ctx context.Context) {
	data := client.CardData{}

	fmt.Print("Введите название карты: ")
	_, err := fmt.Scan(&data.Name)
	if err != nil {
		printErr("Ошибка: не заполнено название карты")
		return
	}

	fmt.Print("Введите номер карты: ")
	_, err = fmt.Scan(&data.Number)
	if err != nil {
		printErr("Ошибка: не заполнен номер карты")
		return
	}

	fmt.Print("Введите месяц окончания действия карты: ")
	_, err = fmt.Scan(&data.Month)
	if err != nil {
		printErr("Ошибка: не заполнен месяц")
		return
	}

	fmt.Print("Введите год окончания действия карты: ")
	_, err = fmt.Scan(&data.Year)
	if err != nil {
		printErr("Ошибка: не заполнен год")
		return
	}

	fmt.Print("Введите CVC карты: ")
	_, err = fmt.Scan(&data.CVC)
	if err != nil {
		printErr("Ошибка: не заполнен CVC")
		return
	}

	fmt.Print("Введите PIN карты: ")
	_, err = fmt.Scan(&data.PIN)
	if err != nil {
		printErr("Ошибка: не заполнен PIN")
		return
	}

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in := bufio.NewReader(os.Stdin)
	data.Meta, _ = in.ReadString('\n')
	if data.Meta == "\n" {
		data.Meta = ""
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	err = nav.Client.AddCard(rCtx, data)
	if err != nil {
		printErr("Ошибка создания карты: " + err.Error())
		return
	}

	printInfo("Карта добавлена")
}

// cardsDelete удаляет карту клиента
func cardsDelete(ctx context.Context) {
	var id string

	fmt.Print("Введите ID карты: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("Ошибка: не заполнен ID карты")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	err = nav.Client.RemoveCard(rCtx, id)
	if err != nil {
		printErr("Ошибка удаления карты: " + err.Error())
		return
	}

	printInfo("Карта удалена")
}

// cardDetail выводит подробную информацию о карте
func cardDetail(ctx context.Context) {
	var id string

	fmt.Print("Введите ID карты: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		printErr("Ошибка: не заполнен ID карты")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := nav.Client.CardDetail(rCtx, id)
	if err != nil {
		printErr("Ошибка получения подробной информации о карте: " + err.Error())
		return
	}

	printInfo("Информация о карте:")
	printInfo(fmt.Sprintf("Наименование: %s", res.Name))
	printInfo(fmt.Sprintf("Номер: %s", res.Number))
	printInfo(fmt.Sprintf("окончание действия: %d/%d", res.Month, res.Year))
	printInfo(fmt.Sprintf("PIN: %s", res.PIN))
	printInfo(fmt.Sprintf("CVC: %s", res.CVC))
	printInfo(fmt.Sprintf("Дополнительная информация: %s", res.Meta))
}

// cardsSearch находит карту клиента по заданным параметрам
func cardsSearch(ctx context.Context) {
	var data client.CardSearchData

	fmt.Print("Введите имя карты для поиска: ")
	_, err := fmt.Scanln(&data.Name)
	if err != nil {
		printErr("Ошибка: не заполнено имя карты")
		return
	}

	md := metadata.Pairs("token", nav.Token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	resp, err := nav.Client.CardSearch(rCtx, data)
	if err != nil {
		printErr("Ошибка поиска карты: " + err.Error())
		return
	}

	if len(resp) == 0 {
		printInfo("У вас нет карт")
		return
	}

	fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
	for _, r := range resp {
		fmt.Println("\033[32m" + fmt.Sprintf("%-*s | %-*s | %-*s | %-*s\033[0m", 38, r.ID, 16, r.Name, 24, r.Number, 24, r.Meta))
	}
}
