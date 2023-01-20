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
	})
}
