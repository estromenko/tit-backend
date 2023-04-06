package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

type whoAmIController struct {
	db          *core.Database
	logger      *core.Logger
	userService *services.UserService
}

func (w *whoAmIController) whoAmI() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, _ := c.Locals("user").(*jwt.Token)

		user, err := w.userService.GetUserByToken(c.UserContext(), token)
		if err != nil {
			w.logger.Err(err).Msg("dashboard user selecting")

			return fiber.ErrInternalServerError
		}

		return c.JSON(user)
	}
}

func NewWhoAmIController(
	db *core.Database,
	conf *core.Config,
	logger *core.Logger,
	userService *services.UserService,
) *fiber.App {
	controller := whoAmIController{db, logger, userService}

	app := fiber.New()

	app.Use(middleware.NewRequireAuth(conf))
	app.Use(middleware.NewIsActive(db, logger))

	app.Get("/", controller.whoAmI())

	return app
}
