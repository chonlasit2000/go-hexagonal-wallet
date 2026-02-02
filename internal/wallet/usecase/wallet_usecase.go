package usecase

import (
	"context"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
)

type walletUsecase struct {
	walletRepo domain.WalletRepository
	txRepo     domain.TransactionRepository
	txManager  domain.TransactionManager
}

func NewWalletUsecase(repo domain.WalletRepository, txRepo domain.TransactionRepository, txManager domain.TransactionManager) domain.WalletUsecase {
	return &walletUsecase{walletRepo: repo, txRepo: txRepo, txManager: txManager}
}

func (u *walletUsecase) GetBalance(ctx context.Context, userID string) (*domain.Wallet, error) {
	// เรียก Repo หาตาม UserID
	wallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	return wallet, nil
}

func (u *walletUsecase) TopUp(ctx context.Context, userID string, amount float64) error {
	// 1. ตรวจสอบยอดเงิน (ต้องมากกว่า 0)
	if amount <= 0 {
		return domain.ErrBadRequest
	}

	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. Update Balance
		wallet, err := u.walletRepo.GetByUserID(txCtx, userID)
		if err != nil {
			return domain.ErrNotFound
		} // ต้องหา WalletID ให้เจอก่อน

		if err := u.walletRepo.UpdateBalance(txCtx, userID, amount); err != nil {
			return err
		}

		// 2. Record History (DEPOSIT)
		history := &domain.Transaction{
			WalletID:    wallet.ID, // ใช้ Wallet ID (UUID)
			Type:        domain.TransactionTypeDeposit,
			Amount:      amount,
			Description: "Top-up via API",
		}
		return u.txRepo.Create(txCtx, history)
	})
}

func (u *walletUsecase) Transfer(ctx context.Context, senderID string, receiverID string, amount float64) error {
	// 1. Validation เบื้องต้น
	if amount <= 0 {
		return domain.ErrBadRequest
	}
	if senderID == receiverID {
		return domain.ErrSameAccount // สร้าง error นี้ใน domain ด้วยนะ
	}

	// 2. เริ่ม Transaction (All or Nothing)
	return u.txManager.Do(ctx, func(txCtx context.Context) error {
		// 1. ดึง Wallet ทั้งคู่ (เพื่อเอา WalletID ที่เป็น UUID)
		senderWallet, err := u.walletRepo.GetByUserID(txCtx, senderID)
		if err != nil {
			return err
		}

		receiverWallet, err := u.walletRepo.GetByUserID(txCtx, receiverID)
		if err != nil {
			return err
		} // จริงๆ ควรเช็คตั้งแต่ validation

		if senderWallet.Balance < amount {
			return domain.ErrInsufficientBalance
		}

		// 2. ปรับยอดเงิน
		if err := u.walletRepo.UpdateBalance(txCtx, senderID, -amount); err != nil {
			return err
		}
		if err := u.walletRepo.UpdateBalance(txCtx, receiverID, amount); err != nil {
			return err
		}

		// 3. บันทึกประวัติฝั่ง "คนโอน" (WITHDRAW)
		logSender := &domain.Transaction{
			WalletID:    senderWallet.ID,
			Type:        domain.TransactionTypeWithdraw,
			Amount:      amount,
			ReferenceID: receiverWallet.ID,
			Description: "Transfer to " + receiverWallet.ID, // หรือใส่ชื่อคนรับถ้ามี
		}
		if err := u.txRepo.Create(txCtx, logSender); err != nil {
			return err
		}

		// 4. บันทึกประวัติฝั่ง "คนรับ" (DEPOSIT)
		logReceiver := &domain.Transaction{
			WalletID:    receiverWallet.ID,
			Type:        domain.TransactionTypeDeposit,
			Amount:      amount,
			ReferenceID: senderWallet.ID,
			Description: "Received from " + senderWallet.ID,
		}
		if err := u.txRepo.Create(txCtx, logReceiver); err != nil {
			return err
		}

		return nil
	})
}

func (u *walletUsecase) GetTransactions(ctx context.Context, userID string) ([]domain.Transaction, error) {
	// 1. หา WalletID ของ user ก่อน
	wallet, err := u.walletRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. ดึง Transaction ของ Wallet นั้น
	return u.txRepo.GetByWalletID(ctx, wallet.ID)
}
