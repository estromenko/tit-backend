package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func NewErrorHandlerMiddleware() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError

		var e *fiber.Error

		if errors.As(err, &e) {
			code = e.Code
		}

		return c.Status(code).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
}
