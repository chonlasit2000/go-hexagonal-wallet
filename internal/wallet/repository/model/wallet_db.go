package model

import (
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/basemodel"
	"github.com/google/uuid"
)

type WalletDBModel struct {
	basemodel.Model           // Embed BaseModel (ได้ UUID, CreatedAt มาฟรีๆ)
	UserID          uuid.UUID `gorm:"type:uuid;not null;index"` // Foreign Key
	Balance         float64   `gorm:"default:0;not null"`
}

func (WalletDBModel) TableName() string {
	return "wallets"
}
