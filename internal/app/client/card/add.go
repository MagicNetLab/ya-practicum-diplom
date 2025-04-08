package card

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
)

// Add добавление новой карты
func addAction(ctx context.Context, cmd *cli.Command, cardClient pb.CardClient) error {
	token, err := readTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	var name, number, month, year, cvc, pin, meta string
	fmt.Print("Введите название карты: ")
	_, err = fmt.Scan(&name)
	if err != nil || name == "" {
		fmt.Println("Ошибка: не заполнено название карты")
		return nil
	}

	fmt.Print("Введите номер карты: ")
	_, err = fmt.Scan(&number)
	if err != nil || number == "" {
		fmt.Println("Ошибка: не заполнен номер карты")
		return nil
	}

	fmt.Print("Введите месяц окончания действия карты: ")
	_, err = fmt.Scan(&month)
	if err != nil {
		fmt.Println("Ошибка: не заполнен месяц")
		return nil
	}
	cardMonth, err := strconv.Atoi(month)
	if err != nil {
		fmt.Println("Ошибка: не заполнен месяц")
		return nil
	}

	fmt.Print("Введите год окончания действия карты: ")
	_, err = fmt.Scan(&year)
	if err != nil {
		fmt.Println("Ошибка: не заполнен год")
		return nil
	}
	cardYear, err := strconv.Atoi(year)
	if err != nil {
		fmt.Println("Ошибка: не заполнен год")
		return nil
	}

	fmt.Print("Введите CVC карты: ")
	_, err = fmt.Scan(&cvc)
	if err != nil {
		fmt.Println("Ошибка: не заполнен CVC")
		return nil
	}

	fmt.Print("Введите PIN карты: ")
	_, err = fmt.Scan(&pin)
	if err != nil {
		fmt.Println("Ошибка: не заполнен PIN")
		return nil
	}

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in := bufio.NewReader(os.Stdin)
	meta, _ = in.ReadString('\n')
	if meta == "\n" {
		meta = ""
	}

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	_, err = cardClient.Create(rCtx, &pb.CreateCardRequest{
		Name:   name,
		Number: number,
		Meta:   meta,
		Month:  int32(cardMonth),
		Year:   int32(cardYear),
		CVC:    cvc,
		PIN:    pin,
	})
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}

		if status.Convert(err).Code() == codes.InvalidArgument {
			fmt.Println("Карта не добавлена. Проверьте правильность заполнения данных.")
			return nil
		}

		fmt.Printf("Ошибка при добавлении карты: %v\n", err)
		return nil
	}

	fmt.Println("Карта успешно добавлена")
	return nil
}
