package http

import (
	"errors"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/response"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	Usecase domain.UserUsecase
}

func NewUserHandler(app *fiber.App, us domain.UserUsecase) {
	handler := &UserHandler{Usecase: us}

	// Route
	api := app.Group("/api")
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
}

func handleError(c *fiber.Ctx, err error) error {
	switch {
	// ถ้า error ตรงกับที่เรานิยามไว้ใน Domain
	case errors.Is(err, domain.ErrBadRequest):
		return response.Error(c, fiber.StatusBadRequest, "Invalid Body", err.Error())
	case errors.Is(err, domain.ErrUnauthorized):
		return response.Error(c, fiber.StatusUnauthorized, "Unauthorized", err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return response.Error(c, fiber.StatusNotFound, "User Not Found", err.Error())
	case errors.Is(err, domain.ErrConflict):
		return response.Error(c, fiber.StatusConflict, "Conflict", err.Error())
	case errors.Is(err, domain.ErrInvalidEmailOrPassword):
		return response.Error(c, fiber.StatusUnauthorized, "Invalid Email or Password", err.Error())
	default:
		return response.Error(c, fiber.StatusInternalServerError, "Internal Server Error", err.Error())
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	// DTO แบบบ้านๆ (Anonymous Struct)
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return handleError(c, domain.ErrBadRequest)
	}

	if err := h.Usecase.Register(c.Context(), req.FirstName, req.LastName, req.Email, req.Password); err != nil {
		return handleError(c, err)
	}

	return response.Success(c, nil)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return handleError(c, domain.ErrBadRequest)
	}

	user, token, err := h.Usecase.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return handleError(c, domain.ErrInvalidEmailOrPassword)
	}

	// สร้าง Response DTO
	res := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}{
		FullName: user.FirstName + " " + user.LastName,
		Email:    user.Email,
		Token:    token,
	}

	return response.Success(c, res)
}
