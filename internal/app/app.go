package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/tutorin-tech/tit-backend/internal/controllers"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

func Run() {
	conf := core.NewConfig()
	log := core.NewLogger(conf)

	db := core.NewDatabase(conf, log)
	defer func() {
		_ = db.Close()
	}()

	log.Info().Msg("Database connection established successfully")

	userService := services.NewUserService(db, log, conf)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(logger.New())

	app.Mount("/auth", controllers.NewAuthController(db, log, userService))
	app.Mount("/api", controllers.NewWhoAmIController(db, conf, log))

	address := fmt.Sprintf(":%d", conf.Port)

	if err := app.Listen(address); err != nil {
		log.Err(err).Msg("Port listening failed")
	}
}
