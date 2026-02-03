package usecase_test

import (
	"context"
	"testing"

	"github.com/chonlasit2000/e-wallet-hexagonal/domain"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/wallet/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock WalletRepository
type mockWalletRepo struct {
	mock.Mock
}

func (m *mockWalletRepo) Create(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockWalletRepo) GetByUserID(ctx context.Context, userID string) (*domain.Wallet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *mockWalletRepo) UpdateBalance(ctx context.Context, userID string, amount float64) error {
	args := m.Called(ctx, userID, amount)
	return args.Error(0)
}

// Mock TransactionRepository
type mockTxRepo struct {
	mock.Mock
}

func (m *mockTxRepo) Create(ctx context.Context, t *domain.Transaction) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *mockTxRepo) GetByWalletID(ctx context.Context, walletID string) ([]domain.Transaction, error) {
	args := m.Called(ctx, walletID)
	return args.Get(0).([]domain.Transaction), args.Error(1)
}

// Mock TransactionManager (สำคัญ! แค่สั่งให้รัน fn เลย ไม่ต้องเปิด DB จริง)
type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	// เรียก fn ทันที โดยส่ง ctx เดิมเข้าไป (Mock ว่า Transaction เริ่มแล้ว)
	return fn(ctx)
}

func TestTransfer_Success(t *testing.T) {
	// 1. Setup
	mockWallet := new(mockWalletRepo)
	mockTx := new(mockTxRepo)
	mockManager := new(mockTxManager)

	u := usecase.NewWalletUsecase(mockWallet, mockTx, mockManager)
	ctx := context.Background()

	senderID := "user-a"
	receiverID := "user-b"
	amount := 100.0

	// Mock Data
	senderWallet := &domain.Wallet{ID: "wallet-a", UserID: senderID, Balance: 1000}
	receiverWallet := &domain.Wallet{ID: "wallet-b", UserID: receiverID, Balance: 0}

	// 2. Expectation (คาดหวังว่าจะเกิดอะไรขึ้นบ้าง)

	// A. ต้องมีการดึง Wallet ผู้โอน
	mockWallet.On("GetByUserID", ctx, senderID).Return(senderWallet, nil)

	// B. ต้องมีการดึง Wallet ผู้รับ
	mockWallet.On("GetByUserID", ctx, receiverID).Return(receiverWallet, nil)

	// C. ต้องมีการตัดเงินผู้โอน (-100)
	mockWallet.On("UpdateBalance", ctx, senderID, -amount).Return(nil)

	// D. ต้องมีการเพิ่มเงินผู้รับ (+100)
	mockWallet.On("UpdateBalance", ctx, receiverID, amount).Return(nil)

	// E. ต้องมีการบันทึก History 2 รอบ (Withdraw, Deposit)
	mockTx.On("Create", ctx, mock.MatchedBy(func(t *domain.Transaction) bool {
		return t.Type == domain.TransactionTypeWithdraw && t.Amount == amount
	})).Return(nil)

	mockTx.On("Create", ctx, mock.MatchedBy(func(t *domain.Transaction) bool {
		return t.Type == domain.TransactionTypeDeposit && t.Amount == amount
	})).Return(nil)

	// 3. Execution (รันจริง)
	err := u.Transfer(ctx, senderID, receiverID, amount)

	// 4. Assertion (ตรวจผลลัพธ์)
	assert.NoError(t, err)           // ต้องไม่มี error
	mockWallet.AssertExpectations(t) // ต้องถูกเรียกครบทุก function
	mockTx.AssertExpectations(t)
}

func TestTransfer_InsufficientBalance(t *testing.T) {
	// 1. Setup
	mockWallet := new(mockWalletRepo)
	mockTx := new(mockTxRepo)
	mockManager := new(mockTxManager) // Mock Manager

	u := usecase.NewWalletUsecase(mockWallet, mockTx, mockManager)
	ctx := context.Background()

	senderID := "user-poor"
	receiverID := "user-rich"
	amount := 5000.0 // โอนเยอะเกิน

	// Mock Data (มีเงินแค่ 100)
	senderWallet := &domain.Wallet{ID: "wallet-poor", UserID: senderID, Balance: 100}

	// 2. Expectation
	// เรียกหาคนโอน -> เจอ
	mockWallet.On("GetByUserID", ctx, senderID).Return(senderWallet, nil)

	// *** ไม่ต้อง Mock ส่วนอื่นต่อ เพราะ Logic ควรจะหยุดตั้งแต่เช็คเงินแล้ว ***

	// 3. Execution
	err := u.Transfer(ctx, senderID, receiverID, amount)

	// 4. Assertion
	assert.Error(t, err)
	assert.Equal(t, domain.ErrInsufficientBalance, err) // Error ต้องตรงเป๊ะ

	// ตรวจสอบว่า UpdateBalance ต้อง **ไม่ถูกเรียก**
	mockWallet.AssertNotCalled(t, "UpdateBalance")
}
