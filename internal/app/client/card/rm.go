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

// removeAction удаляет карту с указанным ID
func removeAction(ctx context.Context, cmd *cli.Command, cardClient pb.CardClient) error {
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
	_, err = cardClient.Delete(rCtx, &pb.DeleteCardRequest{ID: id})
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Println("Ошибка удаления карты: " + err.Error())
		return nil
	}

	fmt.Println("Карта успешно удалена")

	return nil
}
