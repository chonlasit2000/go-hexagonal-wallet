package usecase

import (
	"context"
	"errors"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/pkg/utils"
)

type userUsecase struct {
	userRepo   domain.UserRepository
	walletRepo domain.WalletRepository   // [เพิ่ม] ต้องใช้ WalletRepo
	txManager  domain.TransactionManager // [เพิ่ม] ตัวคุม Transaction
	jwtSecret  string
}

func NewUserUsecase(repo domain.UserRepository, walletRepo domain.WalletRepository, txManager domain.TransactionManager, jwtSecret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:   repo,
		walletRepo: walletRepo,
		txManager:  txManager,
		jwtSecret:  jwtSecret,
	}
}

func (u *userUsecase) Register(ctx context.Context, firstName string, lastName string, email string, password string) error {
	// 1. ตรวจสอบว่ามี Email นี้หรือยัง (Optional)
	_, err := u.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return domain.ErrConflict
	}
	if err != domain.ErrNotFound {
		return err
	}

	// 2. Hash Password *** สำคัญมาก
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	newUser := &domain.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  hashedPwd,
	}

	// *** เริ่ม Transaction ***
	// โค้ดใน fn จะทำงานใน Transaction เดียวกัน
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. สร้าง User (ส่ง txCtx เข้าไป)
		if err := u.userRepo.Store(txCtx, newUser); err != nil {
			return err // ถ้า Error มันจะ Rollback ให้เอง
		}

		// 2. สร้าง Wallet ให้ User คนนี้ทันที (ส่ง txCtx เข้าไป)
		if err := u.walletRepo.Create(txCtx, newUser.ID); err != nil {
			return err // ถ้า Error -> Rollback ทั้ง Wallet และ User
		}

		return nil // ถ้า return nil -> Commit
	})
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
	token, err := utils.GenerateToken(user.ID, u.jwtSecret)
	if err != nil {
		return domain.User{}, "", err
	}

	return *user, token, nil
}
