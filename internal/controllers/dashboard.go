package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
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

		user, err := d.userService.GetUserByToken(c.UserContext(), token)
		if err != nil || user == nil {
			d.logger.Err(err).Msg("dashboard user selecting")

			return fiber.ErrInternalServerError
		}

		dashboardData, err := d.dashboardService.GetDashboard(c.UserContext(), user.ID)
		if err != nil && err != services.ErrDashboardIsStopped {
			d.logger.Err(err).Msg("dashboard port receiving")

			return fiber.ErrInternalServerError
		}

		if dashboardData != nil {
			return c.JSON(dashboardData)
		}

		dashboardData, err = d.dashboardService.StartDashboard(c.UserContext(), user.ID)
		if err != nil {
			d.logger.Err(err).Msg("dashboard starting")

			return fiber.ErrInternalServerError
		}

		return c.JSON(dashboardData)
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
