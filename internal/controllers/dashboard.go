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

		isDashboardRunning := d.dashboardService.IsDashboardRunning(c.UserContext(), user)
		if !isDashboardRunning {
			if err := d.dashboardService.StartDashboard(c.UserContext(), user); err != nil {
				d.logger.Err(err).Msg("dashboard starting")

				return fiber.ErrInternalServerError
			}
		}

		return c.JSON(fiber.Map{
			"id":       user.ID,
			"password": user.DashboardPassword,
		})
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
