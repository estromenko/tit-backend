package controllers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/uptrace/bun"
)

type whoAmIController struct {
	db     *bun.DB
	logger *zerolog.Logger
}

func (w *whoAmIController) whoAmI() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, _ := c.Locals("user").(*jwt.Token)
		claims, _ := token.Claims.(jwt.MapClaims)
		userID, _ := claims["userId"].(float64)

		user := new(models.User)

		err := w.db.NewSelect().
			Model(new(models.User)).
			Where("id = ?", userID).
			Scan(c.UserContext(), user)
		if err != nil {
			return c.SendStatus(http.StatusForbidden)
		}

		return c.JSON(user)
	}
}

func NewWhoAmIController(db *bun.DB, conf *core.Config, logger *zerolog.Logger) *fiber.App {
	controller := whoAmIController{db, logger}

	app := fiber.New()

	app.Use(jwtware.New(jwtware.Config{
		ContextKey:    "user",
		SigningMethod: jwt.SigningMethodHS256.Name,
		SigningKey:    []byte(conf.SecretKey),
	}))

	app.Get("/whoami", controller.whoAmI())

	return app
}
