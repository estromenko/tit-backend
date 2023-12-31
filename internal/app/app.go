package app

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/tutorin-tech/tit-backend/internal/controllers"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
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

	dashboardService, err := services.NewDashboardService(db, conf)
	if err != nil {
		log.Err(err).Msgf("Failed to create dashboard service: %s", err.Error())

		return
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.NewErrorHandlerMiddleware(),
	})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"ping": "pong",
		})
	})

	app.Mount("/auth", controllers.NewAuthController(db, log, userService, conf))
	app.Mount("/api/whoami", controllers.NewWhoAmIController(db, conf, log, userService))
	app.Mount("/api/dashboard", controllers.NewDashboardController(
		db, log, userService, dashboardService, conf,
	))
	app.Mount("/api/tutorials", controllers.NewTutorialsController(db, conf, log, userService))

	address := fmt.Sprintf(":%d", conf.Port)

	if err := app.Listen(address); err != nil {
		log.Err(err).Msg("Port listening failed")
	}
}
