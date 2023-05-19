package main

import (
	"context"
	"fmt"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
	"os"
)

func main() {
	conf := core.NewConfig()
	log := core.NewLogger(conf)

	db := core.NewDatabase(conf, log)
	defer func() {
		_ = db.Close()
	}()

	userService := services.NewUserService(db, log, conf)

	var email, password string

	fmt.Print("Email: ")
	_, _ = fmt.Scanf("%s", &email)

	fmt.Print("Password: ")
	_, _ = fmt.Scanf("%s", &password)

	passwordHash, err := userService.HashPassword(password)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	user := models.User{
		Email:        email,
		IsSuperUser:  true,
		PasswordHash: passwordHash,
	}

	_, err = db.NewInsert().Model(&user).Exec(context.Background())
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
