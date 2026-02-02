package repository

import (
	"context"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/wallet/repository/model"
	"github.com/chonlasit2000/e-wallet-hexagonal/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.WalletRepository {
	return &postgresRepository{db: db}
}

func toDomain(m *model.WalletDBModel) *domain.Wallet {
	return &domain.Wallet{
		ID:      m.Uid.String(),
		UserID:  m.UserID.String(),
		Balance: m.Balance,
	}
}

// func toModel(d *domain.Wallet) *model.WalletDBModel {
// 	userUUID, _ := uuid.Parse(d.UserID)
// 	return &model.WalletDBModel{
// 		Model: basemodel.Model{
// 			Uid: uuid.MustParse(d.ID),
// 		},
// 		UserID:  userUUID,
// 		Balance: d.Balance,
// 	}
// }

func (r *postgresRepository) Create(ctx context.Context, userID string) error {
	userUUID, _ := uuid.Parse(userID)
	walletModel := &model.WalletDBModel{
		UserID:  userUUID,
		Balance: 0,
	}
	db := database.GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Create(walletModel).Error
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
	userUUID, _ := uuid.Parse(userID)
	var dbWallet model.WalletDBModel
	db := database.GetDBFromContext(ctx, r.db)
	if err := db.WithContext(ctx).Where("user_id = ?", userUUID).First(&dbWallet).Error; err != nil {
		return nil, err
	}

	return toDomain(&dbWallet), nil
}

func (r *postgresRepository) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	userUUID, _ := uuid.Parse(userID)
	db := database.GetDBFromContext(ctx, r.db)
	return db.WithContext(ctx).Model(&model.WalletDBModel{}).
		Where("user_id = ?", userUUID).
		Update("balance", gorm.Expr("balance + ?", amount)).Error
}
