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

// searchAction поиск аккаунтов по адресу сайта
func searchAction(ctx context.Context, cmd *cli.Command, accountClient pb.AccountsClient) error {
	var search string
	fmt.Print("Введите Url сайта: ")
	_, err := fmt.Scanln(&search)
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

	req := &pb.SearchAccountRequest{Search: search}
	resp, err := accountClient.Search(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
			return nil
		}
		fmt.Println("Ошибка при получении списка аккаунтов.")
		return nil
	}

	if len(resp.GetAcc()) == 0 {
		fmt.Println("Ни чего не найдено.")
		return nil
	}

	fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, "ID", 16, "Login", 24, "Url", 24, "Desc"))
	for _, row := range resp.GetAcc() {
		fmt.Println(fmt.Sprintf("%-*s | %-*s | %-*s | %-*s", 38, row.GetId(), 16, row.GetLogin(), 24, row.GetUrl(), 24, row.GetDescription()))
	}

	return nil
}
