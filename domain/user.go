package domain

import (
	"context"
	"time"
)

// User Entity: สนใจแค่ Business Logic
type User struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserRepository interface {
	Store(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type UserUsecase interface {
	Register(ctx context.Context, firstName, lastName, email, password string) error
	Login(ctx context.Context, email, password string) (User, string, error)
}
