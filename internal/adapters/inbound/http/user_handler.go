package http

import (
	"net/http"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service inbound.UserService
}

func NewUserHandler(s inbound.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(e *echo.Group) {
	e.POST("/users", h.CreateUser)
	e.GET("/user/:id", h.GetUser)
	e.GET("/users", h.GetUsers)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user domain.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	err := h.service.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id := c.Param("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "user not found"})
	}
	return c.JSON(http.StatusOK, user)
}
func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}
