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

// detailAction выводит подробную информацию о карте с заданным ID.
func detailAction(ctx context.Context, cmd *cli.Command, cardClient pb.CardClient) error {
	token, err := jwt.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}

	id := cmd.Args().Get(0)
	if id == "" {
		fmt.Println("Ошибка: не указан ID карты")
		return nil
	}

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	res, err := cardClient.Get(rCtx, &pb.GetCardRequest{ID: id})
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}

		fmt.Println("Ошибка получения подробной информации о карте: " + err.Error())
		return nil
	}

	card := res.GetCard()

	fmt.Println("Информация о карте:")
	fmt.Println(fmt.Sprintf("Наименование: %s", card.GetName()))
	fmt.Println(fmt.Sprintf("Номер: %s", card.GetNumber()))
	fmt.Println(fmt.Sprintf("Срок действия: %d/%d", card.GetMonth(), card.GetYear()))
	fmt.Println(fmt.Sprintf("PIN: %s", card.GetPIN()))
	fmt.Println(fmt.Sprintf("CVC: %s", card.GetCVC()))
	fmt.Println(fmt.Sprintf("Дополнительная информация: %s", card.GetMeta()))

	return nil
}
