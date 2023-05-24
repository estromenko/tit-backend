package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/tutorin-tech/tit-backend/internal/core"
	"github.com/tutorin-tech/tit-backend/internal/models"
	"github.com/tutorin-tech/tit-backend/internal/services"
	"net/http/httptest"
	"testing"
)

func setupComponents() (
	sqlmock.Sqlmock,
	*services.UserService,
	*fiber.App,
) {
	db, mock := core.NewMockDatabase()
	config := core.NewConfig()
	log := core.NewLogger(config)
	userService := services.NewUserService(db, log, config)
	controller := NewTutorialsController(db, config, log, userService)

	return mock, userService, controller
}

func TestListTutorials(t *testing.T) {
	mock, userService, controller := setupComponents()

	mock.ExpectQuery("SELECT \"u\".\"is_active\"").
		WillReturnRows(sqlmock.NewRows([]string{"is_active"}).AddRow(true))
	mock.ExpectQuery(
		"SELECT \"tuts\".\"id\"",
	).WillReturnRows(sqlmock.NewRows([]string{"id"}))

	req := httptest.NewRequest("GET", "/", nil)
	token, _ := userService.CreateToken(&models.User{})
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := controller.Test(req, 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, fiber.StatusOK)
}

func TestGetTutorial(t *testing.T) {
	mock, userService, controller := setupComponents()

	mock.ExpectQuery("SELECT \"u\".\"is_active\"").
		WillReturnRows(sqlmock.NewRows([]string{"is_active"}).AddRow(true))
	mock.ExpectQuery(
		"SELECT \"tuts\".\"id\"",
	).WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Test tutorial"))

	req := httptest.NewRequest("GET", "/1", nil)
	token, _ := userService.CreateToken(&models.User{})
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	response, err := controller.Test(req, 1)
	assert.Equal(t, err, nil)
	assert.Equal(t, response.StatusCode, fiber.StatusOK)

	tutorial := models.Tutorial{}
	_ = json.NewDecoder(response.Body).Decode(&tutorial)
	assert.Equal(t, tutorial.Name, "Test tutorial")
}
