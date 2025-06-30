package http

import (
	"net/http"

	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/domain"
	"github.com/Aritiaya50217/Backend-Golang-Coding-Test/internal/ports/inbound"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service inbound.AuthenService
}

func NewAuthHandler(s inbound.AuthenService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/login", h.Login)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (h *AuthHandler) CreateUser(c echo.Context) error {
	userID := c.Get("userID").(string)

	authorized, err := h.service.Authorize(userID, "create_user")
	if err != nil || !authorized {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "permission denie"})
	}
	var req domain.User

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "success"})
}
