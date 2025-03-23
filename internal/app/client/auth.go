package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc/status"
)

// auth выполняет авторизацию пользователя
func auth(ctx context.Context) error {
	var login string
	var password string

	fmt.Print("Введите логин: ")
	_, err := fmt.Scanln(&login)
	if err != nil {
		printErr("ОШИБКА: Логин не введен")
		return err
	}
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}

	jwtToken, err := nav.Client.Auth(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			printErr("Ошибка аутентификации: " + e.Message())
			return err
		}
		printErr("Ошибка аутентификации")
		return err
	}

	printSuccess("Аутентификация прошла успешно")
	nav.Token = jwtToken
	clearScreen()
	nav.PrintAvailableCommands()

	return nil
}

// register выполняет регистрацию пользователя
func register(ctx context.Context) error {
	var login string
	var password string
	var rePassword string

	fmt.Print("Введите логин: ")
	_, err := fmt.Scanln(&login)
	if err != nil {
		printErr("Ошибка: логин не введен")
		return err
	}
	fmt.Print("Введите пароль: ")
	_, err = fmt.Scanln(&password)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}
	fmt.Print("Введите пароль еще раз: ")
	_, err = fmt.Scanln(&rePassword)
	if err != nil {
		printErr("ОШИБКА: Пароль не введен")
		return err
	}

	if password != rePassword {
		printErr("ОШИБКА: Пароли не совпадают")
		return fmt.Errorf("пароли не совпадают")
	}

	jwtToken, err := nav.Client.Register(ctx, login, password)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			printErr("Ошибка регистрации: " + e.Message())
			return err
		}
		printErr("Ошибка регистрации")
		return err
	}
	nav.Token = jwtToken

	clearScreen()

	printSuccess("Пользователь успешно зарегистрирован")
	nav.PrintAvailableCommands()

	return nil
}
