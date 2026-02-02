package repository

import (
	"context"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/transaction/repository/model"
	"github.com/chonlasit2000/e-wallet-hexagonal/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.TransactionRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, t *domain.Transaction) error {
	walletUID, _ := uuid.Parse(t.WalletID)

	var refUID *uuid.UUID
	if t.ReferenceID != "" {
		uid, _ := uuid.Parse(t.ReferenceID)
		refUID = &uid
	}

	dbModel := &model.TransactionDBModel{
		WalletID:    walletUID,
		Type:        string(t.Type),
		Amount:      t.Amount,
		ReferenceID: refUID,
		Description: t.Description,
	}

	// ใช้ GetDBFromContext เพื่อรองรับ Transaction
	db := database.GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(dbModel).Error
}

func (r *postgresRepository) GetByWalletID(ctx context.Context, walletID string) ([]domain.Transaction, error) {
	var models []model.TransactionDBModel

	// ดึงรายการล่าสุดขึ้นก่อน
	err := r.db.WithContext(ctx).
		Where("wallet_id = ?", walletID).
		Order("created_at desc").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	// Convert to Domain
	var transactions []domain.Transaction
	for _, m := range models {
		refID := ""
		if m.ReferenceID != nil {
			refID = m.ReferenceID.String()
		}

		transactions = append(transactions, domain.Transaction{
			ID:          m.Uid.String(),
			WalletID:    m.WalletID.String(),
			Type:        domain.TransactionType(m.Type),
			Amount:      m.Amount,
			ReferenceID: refID,
			Description: m.Description,
			CreatedAt:   m.CreatedAt,
		})
	}
	return transactions, nil
}
