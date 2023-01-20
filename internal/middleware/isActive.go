package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/uptrace/bun"
)

func NewIsActive(db *bun.DB, logger *zerolog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, _ := c.Locals("user").(*jwt.Token)
		claims, _ := token.Claims.(jwt.MapClaims)
		userID, _ := claims["userId"].(float64)

		var isActive bool

		err := db.NewSelect().
			Model(new(models.User)).
			Where("id = ?", userID).
			Column("is_active").
			Scan(c.UserContext(), &isActive)
		if err != nil {
			logger.Err(err).Msg("is active middleware")

			return fiber.ErrInternalServerError
		}

		if !isActive {
			return fiber.ErrForbidden
		}

		return c.Next()
	}
}
