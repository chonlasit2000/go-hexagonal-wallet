package http

import (
	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/middleware"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/response"
	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	Usecase domain.WalletUsecase
}

func NewWalletHandler(app *fiber.App, us domain.WalletUsecase, jwtSecret string) {
	handler := &WalletHandler{Usecase: us}

	// สร้าง Group ใหม่ที่ต้องผ่าน Auth Middleware
	api := app.Group("/api/wallets", middleware.Protected(jwtSecret))

	api.Get("/balance", handler.GetBalance)
	api.Post("/topup", handler.TopUp)
	api.Post("/transfer", handler.Transfer)
	api.Get("/history", handler.GetHistory)
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	// 1. ดึง UserID จาก Token (ดู type ดีๆ ว่า Middleware ฝากเป็น string หรือ float64)
	// สมมติใน JWT เราเก็บเป็น string ก็ cast เป็น string
	userID := c.Locals("user_id").(string)

	wallet, err := h.Usecase.GetBalance(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusNotFound, "Wallet not found", err.Error())
	}

	return response.Success(c, wallet)
}

func (h *WalletHandler) TopUp(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid Body", err.Error())
	}

	if err := h.Usecase.TopUp(c.Context(), userID, req.Amount); err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Topup Failed", err.Error())
	}

	return response.Success(c, map[string]string{"status": "success"})
}

func (h *WalletHandler) Transfer(c *fiber.Ctx) error {
	// 1. คนโอนคือคนที่ Login อยู่
	senderID := c.Locals("user_id").(string)

	// 2. รับข้อมูลคนรับ + จำนวนเงิน
	var req struct {
		ReceiverID string  `json:"receiver_id" validate:"required"`
		Amount     float64 `json:"amount" validate:"required,min=1"` // ห้ามโอน 0 หรือติดลบ
	}

	if err := c.BodyParser(&req); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid Body", err.Error())
	}

	// (ควรเรียก utils.ValidateStruct ตรงนี้ด้วย ถ้าทำ validation แล้ว)

	// 3. เรียก Usecase
	err := h.Usecase.Transfer(c.Context(), senderID, req.ReceiverID, req.Amount)
	if err != nil {
		// Map Error ให้ตรงกับ HTTP Status
		switch err {
		case domain.ErrInsufficientBalance:
			return response.Error(c, fiber.StatusBadRequest, "Insufficient Balance", err.Error())
		case domain.ErrNotFound:
			return response.Error(c, fiber.StatusNotFound, "Receiver Not Found", err.Error())
		case domain.ErrSameAccount:
			return response.Error(c, fiber.StatusBadRequest, "Cannot Transfer to Self", err.Error())
		default:
			return response.Error(c, fiber.StatusInternalServerError, "Transfer Failed", err.Error())
		}
	}

	return response.Success(c, map[string]string{"status": "transfer success"})
}

func (h *WalletHandler) GetHistory(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	txs, err := h.Usecase.GetTransactions(c.Context(), userID)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Failed to get history", err.Error())
	}

	return response.Success(c, txs)
}
