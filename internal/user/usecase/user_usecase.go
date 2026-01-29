package usecase

import (
	"context"
	"errors"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/pkg/utils"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) Register(ctx context.Context, firstName string, lastName string, email string, password string) error {
	// 1. ตรวจสอบว่ามี Email นี้หรือยัง (Optional)
	existingUser, _ := u.userRepo.GetByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// 2. Hash Password *** สำคัญมาก
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	// 3. ส่งให้ Repo บันทึก
	newUser := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPwd,
	}

	return u.userRepo.Store(ctx, newUser)
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (domain.User, string, error) {
	// 1. หา User
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return domain.User{}, "", errors.New("invalid email or password")
	}

	// 2. เช็ค Password
	if !utils.CheckPasswordHash(password, user.Password) {
		return domain.User{}, "", errors.New("invalid email or password")
	}

	// 3. แจก Token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return domain.User{}, "", err
	}

	return *user, token, nil
}
