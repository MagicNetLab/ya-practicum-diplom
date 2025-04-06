package card

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/card/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// searchAction получает список карт пользователя и выводит его в консоль
func searchAction(ctx context.Context, cmd *cli.Command, cardClient pb.CardClient) error {
	token, err := jwt.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	var search string
	fmt.Print("Введите имя карты для поиска: ")
	_, err = fmt.Scanln(&search)
	if err != nil {
		fmt.Println("Ошибка: не заполнено имя карты")
		return nil
	}

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	req := &pb.SearchCardRequest{Name: search}
	res, err := cardClient.Search(rCtx, req)

	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Println("Ошибка получения списка карт: " + err.Error())
		return nil
	}

	if len(res.GetCards()) == 0 {
		fmt.Println("Ничего не найдено")
		return nil
	}

	fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
	for _, r := range res.GetCards() {
		fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, r.GetID(), 16, r.GetName(), 24, r.GetNumber(), 24, r.GetMeta()))
	}

	return nil
}
