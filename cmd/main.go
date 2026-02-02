package main

import (
	"fmt"
	"log"

	"github.com/chonlasit2000/e-wallet-hexagonal/config"
	transactionRepo "github.com/chonlasit2000/e-wallet-hexagonal/internal/transaction/repository"
	modeltransaction "github.com/chonlasit2000/e-wallet-hexagonal/internal/transaction/repository/model"
	userHttp "github.com/chonlasit2000/e-wallet-hexagonal/internal/user/delivery/http"
	userRepo "github.com/chonlasit2000/e-wallet-hexagonal/internal/user/repository"
	userUsecase "github.com/chonlasit2000/e-wallet-hexagonal/internal/user/usecase"
	walletHttp "github.com/chonlasit2000/e-wallet-hexagonal/internal/wallet/delivery/http"
	walletRepo "github.com/chonlasit2000/e-wallet-hexagonal/internal/wallet/repository"
	walletUsecase "github.com/chonlasit2000/e-wallet-hexagonal/internal/wallet/usecase"
	"github.com/chonlasit2000/e-wallet-hexagonal/pkg/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("✅ Config loaded")

	// 2. Connect Database
	db := database.NewPostgresDatabase(cfg)
	fmt.Println("✅ Database connected")

	// Auto Migrate (สร้าง Table ให้อัตโนมัติ)
	// db.AutoMigrate(&modeluser.UserDBModel{})
	// db.AutoMigrate(&modelwallet.WalletDBModel{})
	db.AutoMigrate(&modeltransaction.TransactionDBModel{})
	fmt.Println("✅ Database migrated")

	// 3. Init Repositories
	userRepo := userRepo.NewPostgresRepository(db)
	walletRepo := walletRepo.NewPostgresRepository(db)
	transactionRepo := transactionRepo.NewPostgresRepository(db)

	// Init Transactor
	txManager := database.NewGormTransactionManager(db)

	// Init Usecases
	userUsecase := userUsecase.NewUserUsecase(userRepo, walletRepo, txManager, cfg.Server.JWTSecret)
	walletUsecase := walletUsecase.NewWalletUsecase(walletRepo, transactionRepo, txManager)

	// 4. Start Server
	app := fiber.New()

	userHttp.NewUserHandler(app, userUsecase)
	walletHttp.NewWalletHandler(app, walletUsecase, cfg.Server.JWTSecret)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("E-Wallet API is Running! 🚀")
	})

	// ใช้ Port จาก Config
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
