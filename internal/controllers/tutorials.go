package controllers

import (
	"database/sql"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/middleware"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
)

type tutorialsController struct {
	db          *core.Database
	logger      *core.Logger
	userService *services.UserService
}

func (t *tutorialsController) listTutorials() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var tutorials []*models.Tutorial

		_, err := t.db.NewSelect().Model(new(models.Tutorial)).Exec(c.UserContext(), &tutorials)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.logger.Err(err).Msg("list tutorials")

			return fiber.ErrInternalServerError
		}

		if tutorials == nil {
			tutorials = []*models.Tutorial{}
		}

		return c.JSON(tutorials)
	}
}

func (t *tutorialsController) createTutorial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tutorial := new(models.Tutorial)
		if err := c.BodyParser(&tutorial); err != nil {
			return err
		}

		_, err := t.db.NewSelect().
			Model(tutorial).
			Where("name = ?", tutorial.Name).
			Exec(c.UserContext(), tutorial)
		if err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "tutorial with such name already exists",
			})
		}

		if !errors.Is(err, sql.ErrNoRows) {
			t.logger.Err(err).Msg("existent tutorial selecting")

			return fiber.ErrInternalServerError
		}

		_, err = t.db.NewInsert().
			Model(tutorial).
			Returning("id").
			Exec(c.UserContext(), &tutorial.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			t.logger.Err(err).Msg("tutorial insert")

			return fiber.ErrInternalServerError
		}

		return c.JSON(tutorial)
	}
}

func (t *tutorialsController) getTutorial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tutorial := new(models.Tutorial)

		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrNotFound
		}

		tutorial.ID = uint64(id)

		err = t.db.NewSelect().Model(tutorial).WherePK().Scan(c.UserContext(), tutorial)
		if errors.Is(err, sql.ErrNoRows) {
			return fiber.ErrNotFound
		}

		return c.JSON(tutorial)
	}
}

func (t *tutorialsController) updateTutorial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tutorial := new(models.Tutorial)

		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrNotFound
		}

		if err := c.BodyParser(tutorial); err != nil {
			return err
		}

		tutorial.ID = uint64(id)

		existentTutorial := new(models.Tutorial)
		_, err = t.db.NewSelect().Model(tutorial).Where("id = ?", id).Exec(c.UserContext(), existentTutorial)
		if err != nil {
			if err == sql.ErrNoRows {
				return fiber.ErrNotFound
			}

			t.logger.Err(err).Msg("existent tutorial selecting")

			return fiber.ErrInternalServerError
		}

		_, err = t.db.NewUpdate().Model(tutorial).WherePK().Exec(c.UserContext())
		if err != nil {
			t.logger.Err(err).Msg("tutorial update")

			return fiber.ErrInternalServerError
		}

		return c.JSON(tutorial)
	}
}

func (t *tutorialsController) deleteTutorial() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tutorial := new(models.Tutorial)

		id, err := c.ParamsInt("id")
		if err != nil {
			return fiber.ErrNotFound
		}

		tutorial.ID = uint64(id)

		_, err = t.db.NewDelete().Model(tutorial).WherePK().Exec(c.UserContext())
		if err != nil {
			t.logger.Err(err).Msg("tutorial delete")

			return fiber.ErrInternalServerError
		}

		return c.JSON(fiber.Map{
			"message": "tutorial deleted successfully",
		})
	}
}

func NewTutorialsController(
	db *core.Database,
	conf *core.Config,
	logger *core.Logger,
	userService *services.UserService,
) *fiber.App {
	controller := tutorialsController{db, logger, userService}

	app := fiber.New()

	app.Use(middleware.NewRequireAuth(conf))
	app.Use(middleware.NewIsActive(db, logger))

	app.Get("/", controller.listTutorials())
	app.Post("/", controller.createTutorial())
	app.Get("/:id", controller.getTutorial())
	app.Put("/:id", controller.updateTutorial())
	app.Delete("/:id", controller.deleteTutorial())

	return app
}
