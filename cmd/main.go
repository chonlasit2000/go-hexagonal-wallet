package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chonlasit2000/e-wallet-hexagonal/config"
	transactionRepo "github.com/chonlasit2000/e-wallet-hexagonal/internal/transaction/repository"
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
	// Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	fmt.Println("✅ Config loaded")

	// Connect Database
	db := database.NewPostgresDatabase(cfg)
	fmt.Println("✅ Database connected")

	// Auto Migrate (สร้าง Table ให้อัตโนมัติ)
	// db.AutoMigrate(&modeluser.UserDBModel{})
	// db.AutoMigrate(&modelwallet.WalletDBModel{})
	// db.AutoMigrate(&modeltransaction.TransactionDBModel{})
	fmt.Println("✅ Database migrated")

	// Init Repositories
	userRepo := userRepo.NewPostgresRepository(db)
	walletRepo := walletRepo.NewPostgresRepository(db)
	transactionRepo := transactionRepo.NewPostgresRepository(db)

	// Init Transactor
	txManager := database.NewGormTransactionManager(db)

	// Init Usecases
	userUsecase := userUsecase.NewUserUsecase(userRepo, walletRepo, txManager, cfg.Server.JWTSecret)
	walletUsecase := walletUsecase.NewWalletUsecase(walletRepo, transactionRepo, txManager)

	app := fiber.New()
	userHttp.NewUserHandler(app, userUsecase)
	walletHttp.NewWalletHandler(app, walletUsecase, cfg.Server.JWTSecret)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("E-Wallet API is Running! 🚀")
	})

	// --- Graceful Shutdown Start ---
	// สร้าง Channel รอรับสัญญาณ
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // ดัก Ctrl+C

	// รัน Server ใน Goroutine (แยกไปทำงานเบื้องหลัง ไม่ให้ Block Main Thread)
	go func() {
		fmt.Printf("Server is running on port %s 🚀\n", cfg.Server.Port)
		if err := app.Listen(":" + cfg.Server.Port); err != nil {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	// Main Thread จะมารอตรงนี้ (Block จนกว่าจะมีสัญญาณ)
	<-quit
	fmt.Println("\nShutting down server...")

	// สั่งปิด Server (ให้เวลาเคลียร์ Request 5 วินาที)
	// (ถ้า Fiber เวอร์ชั่นใหม่ๆ อาจต้องใช้ context.WithTimeout แต่นี่ใช้แบบง่ายไปก่อนได้ครับ)
	if err := app.Shutdown(); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	fmt.Println("Server exiting")
}
