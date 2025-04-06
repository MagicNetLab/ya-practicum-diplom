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

// listAction получает список карт и выводит его в консоль.
func listAction(ctx context.Context, cmd *cli.Command, cardClient pb.CardClient) error {
	token, err := jwt.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return nil
	}
	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	req := &pb.ListCardRequest{}
	resp, err := cardClient.List(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}

		fmt.Printf("Ошибка получения списка карт: %v\n", err)
		return nil
	}

	if len(resp.GetCards()) == 0 {
		fmt.Println("У вас нет карт")
		return nil
	}

	fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Name", 24, "Number", 24, "Desc"))
	for _, row := range resp.GetCards() {
		fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, row.GetID(), 16, row.GetName(), 24, row.GetNumber(), 24, row.GetMeta()))
	}

	return nil
}
