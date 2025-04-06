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
	user := cmd.String("user")
	password := cmd.String("password")

	if user == "" || password == "" {
		fmt.Println("не указан пользователь или пароль")
		return fmt.Errorf("не указан пользователь или пароль")
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
