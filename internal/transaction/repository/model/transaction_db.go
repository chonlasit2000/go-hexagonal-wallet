package model

import (
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/basemodel"
	"github.com/google/uuid"
)

type TransactionDBModel struct {
	basemodel.Model
	WalletID    uuid.UUID  `gorm:"type:uuid;index;not null"` // ของ Wallet ไหน
	Type        string     `gorm:"type:varchar(20);not null"`
	Amount      float64    `gorm:"not null"`
	ReferenceID *uuid.UUID `gorm:"type:uuid;index"` // Wallet คู่กรณี (อาจไม่มีถ้า Topup)
	Description string     `gorm:"type:varchar(255)"`
}

func (TransactionDBModel) TableName() string {
	return "transactions"
}
