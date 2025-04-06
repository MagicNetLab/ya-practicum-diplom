package accounts

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/account/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// addAction добавления нового аккаунта
func addAction(ctx context.Context, cmd *cli.Command, accountClient pb.AccountsClient) error {
	token, err := jwt.ReadTokenFromFile()
	if err != nil {
		fmt.Println("Ошибка при получении токена. Возможно, вы не авторизовались")
		return err
	}

	var login, password, url, meta string

	fmt.Print("Введите Url сайта: ")
	_, err = fmt.Scanln(&url)
	if err != nil {
		fmt.Println("ОШИБКА: Адрес сайта не введен.")
		return err
	}

	fmt.Print("Введите Логин для сайта: ")
	_, err = fmt.Scanln(&login)
	if err != nil {
		fmt.Println("ОШИБКА: Логин не введен.")
		return nil
	}

	fmt.Print("Введите Пароль для сайта: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		fmt.Println("ОШИБКА: Логин не введен.")
		return nil
	}

	fmt.Print("Введите дополнительную информацию (не обязательно): ")
	in := bufio.NewReader(os.Stdin)
	meta, _ = in.ReadString('\n')
	meta = strings.TrimSpace(meta)

	md := metadata.Pairs("token", token)
	rCtx := metadata.NewOutgoingContext(ctx, md)

	req := &pb.CreateAccountRequest{
		Login:       login,
		Password:    password,
		Url:         url,
		Description: meta,
	}

	_, err = accountClient.Create(rCtx, req)
	if err != nil {
		if status.Convert(err).Code() == codes.Unauthenticated {
			_ = os.Remove("token.txt")
			fmt.Println("Время жизни токена истекло. Пожалуйста, повторите авторизацию.")
		}
		return nil
	}

	fmt.Println("Аккаунт успешно добавлен!")

	return nil
}
