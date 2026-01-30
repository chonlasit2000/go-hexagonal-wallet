package domain

import "context"

type Wallet struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

// Repository Interface
type WalletRepository interface {
	Create(ctx context.Context, userID string) error
	GetByUserID(ctx context.Context, userID string) (*Wallet, error)
	UpdateBalance(ctx context.Context, userID string, amount float64) error // ใช้สำหรับเติม/ตัดเงิน
}

// Usecase Interface (สำหรับเรียกใช้)
type WalletUsecase interface {
	GetBalance(ctx context.Context, userID string) (*Wallet, error)
	TopUp(ctx context.Context, userID string, amount float64) error
}
