package repository

import (
	"context"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/basemodel"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/user/repository/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) domain.UserRepository {
	return &postgresRepository{db: db}
}

// ฟังก์ชันแปลง DB Model -> Domain
func toDomain(m *model.UserDBModel) *domain.User {
	return &domain.User{
		ID:        m.Uid.String(), // แปลง uuid.UUID -> string
		FirstName: m.FirstName,
		LastName:  m.LastName,
		Email:     m.Email,
		Password:  m.Password,
		CreatedAt: m.CreatedAt,
	}
}

// ฟังก์ชันแปลง Domain -> DB Model
func toModel(d *domain.User) *model.UserDBModel {
	uid, _ := uuid.Parse(d.ID) // แปลง string กลับเป็น uuid (ถ้ามีค่า)

	return &model.UserDBModel{
		Model: basemodel.Model{
			Uid: uid,
		},
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Email:     d.Email,
		Password:  d.Password,
	}
}

func (r *postgresRepository) Store(ctx context.Context, u *domain.User) error {
	// 1. แปลง Domain เป็น Model เพื่อเตรียมบันทึก
	dbUser := toModel(u)

	// 2. บันทึกลง Postgres
	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return err
	}

	// 3. (Optional) อัปเดต ID กลับไปให้ Domain (เพราะ Postgres เพิ่ง gen ID มา)
	u.ID = dbUser.Uid.String()
	return nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var dbUser model.UserDBModel
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomain(&dbUser), nil
}
