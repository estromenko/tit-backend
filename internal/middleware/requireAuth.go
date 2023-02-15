package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tutorin-tech/tit-backend/internal/core"
)

func NewRequireAuth(conf *core.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		ContextKey:    "user",
		SigningMethod: jwt.SigningMethodHS256.Name,
		SigningKey:    []byte(conf.SecretKey),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Missing or malformed JWT",
				})
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired JWT",
			})
		},
	})
}
