package controllers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
	"github.com/uptrace/bun"
)

type authController struct {
	db          *bun.DB
	logger      *zerolog.Logger
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
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
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
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{
					"error": "Wrong email or password",
				})
			}

			a.logger.Err(err).Msg("login user selecting")

			return c.SendStatus(http.StatusInternalServerError)
		}

		if !user.IsActive {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}

		passwordsMatch, err := a.userService.CheckPassword(user, requestData.Password)
		if err != nil {
			a.logger.Err(err).Msg("login password check")

			return c.SendStatus(http.StatusInternalServerError)
		}

		if !passwordsMatch {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Wrong email or password",
			})
		}

		token, err := a.userService.CreateToken(user)
		if err != nil {
			a.logger.Err(err).Msg("login token creation")

			return c.SendStatus(http.StatusInternalServerError)
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
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if a.userService.UserExists(c.UserContext(), user.Email) {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "User with such email already exists",
			})
		}

		passwordHash, err := a.userService.HashPassword(user.Password)
		if err != nil {
			a.logger.Err(err).Msg("registration password hash")

			return c.SendStatus(http.StatusInternalServerError)
		}

		user.PasswordHash = passwordHash

		tx, _ := a.db.BeginTx(c.UserContext(), &sql.TxOptions{Isolation: 0, ReadOnly: false})

		_, err = tx.NewInsert().Model(user).Exec(c.UserContext())
		if err != nil {
			a.logger.Err(err).Msg("registration user insert")

			return c.SendStatus(http.StatusInternalServerError)
		}

		token, err := a.userService.CreateToken(user)
		if err != nil {
			_ = tx.Rollback()

			a.logger.Err(err).Msg("registration token creation")

			return c.SendStatus(http.StatusInternalServerError)
		}

		if err := tx.Commit(); err != nil {
			a.logger.Err(err).Msg("registration transaction commit")

			return c.SendStatus(http.StatusInternalServerError)
		}

		user.Password = ""
		user.Token = token

		return c.JSON(user)
	}
}

func NewAuthController(
	db *bun.DB,
	logger *zerolog.Logger,
	userService *services.UserService,
) *fiber.App {
	controller := authController{db, logger, userService}

	app := fiber.New()
	app.Post("/login", controller.login())
	app.Post("/registration", controller.registration())

	return app
}
