package controllers

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

type dashboardController struct {
	db               *core.Database
	logger           *core.Logger
	userService      *services.UserService
	dashboardService *services.DashboardService
}

func (d *dashboardController) dashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, _ := c.Locals("user").(*jwt.Token)
		claims, _ := token.Claims.(jwt.MapClaims)
		userID, _ := claims["userId"].(float64)

		user := new(models.User)

		err := d.db.NewSelect().
			Model(new(models.User)).
			Where("id = ?", userID).
			Scan(c.UserContext(), user)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fiber.ErrForbidden
			}

			d.logger.Err(err).Msg("whoami user selecting")

			return fiber.ErrInternalServerError
		}

		port, err := d.dashboardService.GetDashboardPort(c.UserContext(), user.ID)
		if err != nil {
			d.logger.Err(err).Msg("dashboard port receiving")

			return fiber.ErrInternalServerError
		}

		if port != 0 {
			return c.JSON(fiber.Map{"port": port})
		}

		port, err = d.dashboardService.StartDashboard(c.UserContext(), user.ID)
		if err != nil {
			d.logger.Err(err).Msg("dashboard starting")

			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{"port": port})
	}
}

func NewDashboardController(
	db *core.Database,
	logger *core.Logger,
	userService *services.UserService,
	dashboardService *services.DashboardService,
	conf *core.Config,
) *fiber.App {
	controller := dashboardController{
		db,
		logger,
		userService,
		dashboardService,
	}

	app := fiber.New()

	app.Use(middleware.NewRequireAuth(conf))
	app.Use(middleware.NewIsActive(db, logger))

	app.Post("/", controller.dashboard())

	return app
}
