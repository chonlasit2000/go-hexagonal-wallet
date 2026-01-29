package model

// UserDBModel สืบทอดคุณสมบัติจาก BaseModel
type UserDBModel struct {
	BaseModel        // ฝัง struct (Embedding)
	FirstName string `gorm:"type:varchar(100)"`
	LastName  string `gorm:"type:varchar(100)"`
	Email     string `gorm:"uniqueIndex;type:varchar(100)"`
	Password  string `gorm:"type:varchar(255)"`
}

func (UserDBModel) TableName() string {
	return "users"
}
