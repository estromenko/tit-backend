package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/tutorin-tech/tit-backend/internal/core"
)

func Run() {
	conf := core.NewConfig()
	log := core.NewLogger(conf)
	app := fiber.New()

	db := core.NewDatabase(conf)
	defer func() {
		_ = db.Close()
	}()

	if err := db.Ping(); err != nil {
		log.Err(err).Msg("Database ping failed")

		return
	}

	log.Info().Msg("Database connection established successfully")

	address := fmt.Sprintf(":%d", conf.Port)

	if err := app.Listen(address); err != nil {
		log.Err(err).Msg("Port listening failed")
	}
}
