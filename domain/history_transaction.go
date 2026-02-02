package domain

import (
	"context"
	"time"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "DEPOSIT"  // เงินเข้า (เติมเงิน, รับโอน)
	TransactionTypeWithdraw TransactionType = "WITHDRAW" // เงินออก (โอนออก)
)

type Transaction struct {
	ID          string          `json:"id"`
	WalletID    string          `json:"wallet_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	ReferenceID string          `json:"reference_id,omitempty"` // อ้างอิง Wallet คู่กรณี (ถ้ามี)
	Description string          `json:"description"`            // คำอธิบาย (เช่น "Topup", "Transfer to ...")
	CreatedAt   time.Time       `json:"created_at"`
}

type TransactionRepository interface {
	Create(ctx context.Context, t *Transaction) error
	GetByWalletID(ctx context.Context, walletID string) ([]Transaction, error) // ของจริงควรมี Pagination
}
