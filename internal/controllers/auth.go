package controllers

import (
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

type authController struct {
	db          *core.Database
	logger      *core.Logger
	userService *services.UserService
}

func (a *authController) login() fiber.Handler { //nolint:funlen
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c *fiber.Ctx) error {
		requestData := new(request)
		_ = c.BodyParser(requestData)

		validate := validator.New()
		if err := validate.Struct(requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		user := new(models.User)

		err := a.db.NewSelect().
			Model(user).
			Where("email = ?", requestData.Email).
			Scan(c.UserContext(), user)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Wrong email or password",
				})
			}

			a.logger.Err(err).Msg("login user selecting")

			return fiber.ErrInternalServerError
		}

		if !user.IsActive {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}

		passwordsMatch, err := a.userService.CheckPassword(user, requestData.Password)
		if err != nil {
			a.logger.Err(err).Msg("login password check")

			return fiber.ErrInternalServerError
		}

		if !passwordsMatch {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}

		token, err := a.userService.CreateToken(user)
		if err != nil {
			a.logger.Err(err).Msg("login token creation")

			return fiber.ErrInternalServerError
		}

		user.Password = ""
		user.Token = token

		return c.JSON(user)
	}
}

func (a *authController) registration() fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		_ = c.BodyParser(user)
		user.IsActive = true

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if a.userService.UserExists(c.UserContext(), user.Email) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User with such email already exists",
			})
		}

		passwordHash, err := a.userService.HashPassword(user.Password)
		if err != nil {
			a.logger.Err(err).Msg("registration password hash")

			return fiber.ErrInternalServerError
		}

		user.PasswordHash = passwordHash

		tx, _ := a.db.BeginTx(c.UserContext(), &sql.TxOptions{Isolation: 0, ReadOnly: false})

		_, err = tx.NewInsert().Model(user).Exec(c.UserContext())
		if err != nil {
			a.logger.Err(err).Msg("registration user insert")

			return fiber.ErrInternalServerError
		}

		token, err := a.userService.CreateToken(user)
		if err != nil {
			_ = tx.Rollback()

			a.logger.Err(err).Msg("registration token creation")

			return fiber.ErrInternalServerError
		}

		if err := tx.Commit(); err != nil {
			a.logger.Err(err).Msg("registration transaction commit")

			return fiber.ErrInternalServerError
		}

		user.Password = ""
		user.Token = token

		return c.JSON(user)
	}
}

func (a *authController) resetPassword() fiber.Handler {
	type request struct {
		Password    string `json:"password" validate:"required"`
		NewPassword string `json:"newPassword" validate:"required,min=8,max=256"`
	}

	return func(c *fiber.Ctx) error {
		requestData := new(request)
		_ = c.BodyParser(requestData)

		validate := validator.New()
		if err := validate.Struct(requestData); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		token, _ := c.Locals("user").(*jwt.Token)

		user, err := a.userService.GetUserByToken(c.UserContext(), token)
		if err != nil || user == nil {
			a.logger.Err(err).Msg("dashboard user selecting")

			return fiber.ErrInternalServerError
		}

		hashedPassword, err := a.userService.HashPassword(requestData.Password)
		if err != nil {
			a.logger.Err(err).Msg("reset password hashing")

			return fiber.ErrInternalServerError
		}

		if user.PasswordHash != hashedPassword {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "wrong password",
			})
		}

		newHashedPassword, err := a.userService.HashPassword(requestData.NewPassword)
		if err != nil {
			a.logger.Err(err).Msg("reset new password hashing")

			return fiber.ErrInternalServerError
		}

		_, err = a.db.NewUpdate().
			Model(user).
			Set("password_hash = ?", newHashedPassword).
			WherePK().
			Exec(c.UserContext())
		if err != nil {
			a.logger.Err(err).Msg("reset password saving")

			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "password changed successfully",
		})
	}
}

func NewAuthController(
	db *core.Database,
	logger *core.Logger,
	userService *services.UserService,
	conf *core.Config,
) *fiber.App {
	controller := authController{db, logger, userService}

	app := fiber.New()
	app.Post("/login", controller.login())
	app.Post("/registration", controller.registration())

	app.Use(middleware.NewRequireAuth(conf))
	app.Use(middleware.NewIsActive(db, logger))
	app.Post("/reset-password", controller.resetPassword())

	return app
}
