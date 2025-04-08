package accounts

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// detailAction выводит информацию об аккаунте с указанным ID.
func detailAction(ctx context.Context, cmd *cli.Command, accountClient pb.AccountsClient) error {
	var id string
	fmt.Print("Введите Url сайта: ")
	_, err := fmt.Scanln(&id)
	if err != nil {
		fmt.Println("Необходимо указать id аккаунта")
		return err
	}

	token, err := jwt.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return err
	}

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)
	req := &pb.GetAccountRequest{Id: id}

	resp, err := accountClient.Get(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Println("Ошибка при получении информации об аккаунте.")
		return nil
	}

	acc := resp.GetAccount()
	fmt.Println("Данные аккаунта:")
	fmt.Println("Url: " + acc.GetUrl())
	fmt.Println("Логин: " + acc.GetLogin())
	fmt.Println("Пароль: " + acc.GetPassword())
	fmt.Println("Meta: " + acc.GetDescription())

	return nil
}
