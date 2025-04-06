package auth

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"

	pb "github.com/MagicNetLab/ya-practicum-diplom/internal/grpc/auth/proto"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/jwt"
)

// authAction выполняет аутентификацию пользователя и сохраняет токен в файл
func authAction(ctx context.Context, cmd *cli.Command, authClient pb.AuthClient) error {
	user := cmd.String("user")
	password := cmd.String("password")

	if user == "" || password == "" {
		fmt.Println("не указан пользователь или пароль")
		return fmt.Errorf("не указан пользователь или пароль")
	}

	request := &pb.AuthRequest{
		Login:  user,
		Secret: password,
	}

	resp, err := authClient.Auth(ctx, request)
	if err != nil {
		fmt.Printf("ошибка аутентификации: %v", err)
		return fmt.Errorf("ошибка аутентификации: %v", err)
	}

	err = jwt.SaveTokenToFile(resp.GetToken())
	if err != nil {
		return fmt.Errorf("ошибка записи токена: %v", err.Error())
	}
	fmt.Println("Авторизация прошла успешно")

	return nil
}
