package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tutorin-tech/tit-backend/internal/controllers"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

func Run() {
	conf := core.NewConfig()
	log := core.NewLogger(conf)

	db := core.NewDatabase(conf)
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(); err != nil {
		log.Err(err).Msg("Database ping failed")

		return
	}

	log.Info().Msg("Database connection established successfully")

	userService := services.NewUserService(db, log, conf)

	app := fiber.New()
	app.Mount("/auth", controllers.NewAuthController(db, log, userService))

	address := fmt.Sprintf(":%d", conf.Port)

	if err := app.Listen(address); err != nil {
		log.Err(err).Msg("Port listening failed")
	}
}