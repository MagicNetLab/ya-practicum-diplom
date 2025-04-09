package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// RegisterAction - регистрация пользователя
func registerAction(ctx context.Context, cmd *cli.Command, authClient pb.AuthClient) error {
	var user, password string

	fmt.Print("Введите имя пользователя: ")
	_, err := fmt.Scan(&user)
	if err != nil || user == "" {
		fmt.Println("Ошибка: не указанно имя пользователя")
		return nil
	}

	fmt.Print("Введите пароль: ")
	_, err = fmt.Scan(&password)
	if err != nil || password == "" {
		fmt.Println("Ошибка: не указан пароль")
		return nil
	}

	request := &pb.RegRequest{
		Login:  user,
		Secret: password,
	}

	resp, err := authClient.Register(ctx, request)
	if err != nil {
		fmt.Printf("ошибка регистрации: %v", err)
		return err
	}

	err = jwt.SaveTokenToFile(resp.GetToken())
	if err != nil {
		return fmt.Errorf("ошибка записи токена: %v", err.Error())
	}
	fmt.Println("Регистрация прошла успешно")

	return nil
}
