package database

import (
	"context"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"gorm.io/gorm"
)

type TxKey struct{}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) domain.TransactionManager {
	return &gormTransactionManager{db: db}
}

func (t *gormTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		// เอา tx (DB ที่ติด Transaction) ยัดใส่ Context
		txCtx := context.WithValue(ctx, TxKey{}, tx)
		return fn(txCtx)
	})
}

// Helper ให้ Repo เรียกใช้: ถ้ามี tx ใน context ให้ใช้ tx, ถ้าไม่มีให้ใช้ db เดิม
func GetDBFromContext(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(TxKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
