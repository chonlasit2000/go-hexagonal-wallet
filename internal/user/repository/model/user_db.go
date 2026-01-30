package model

import "github.com/chonlasit2000/e-wallet-hexagonal/internal/basemodel"

type UserDBModel struct {
	basemodel.Model        // ฝัง struct (Embedding)
	FirstName       string `gorm:"type:varchar(100)"`
	LastName        string `gorm:"type:varchar(100)"`
	Email           string `gorm:"uniqueIndex;type:varchar(100)"`
	Password        string `gorm:"type:varchar(255)"`
}

func (UserDBModel) TableName() string {
	return "users"
}
