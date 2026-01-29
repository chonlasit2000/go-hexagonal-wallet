package main

import (
	"fmt"
	"log"

	"github.com/chonlasit2000/e-wallet-hexagonal/config"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/user/delivery/http"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/user/repository"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/user/repository/model"
	"github.com/chonlasit2000/e-wallet-hexagonal/internal/user/usecase"
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
	db.AutoMigrate(&model.UserDBModel{})
	fmt.Println("✅ Database migrated")

	// 3. Init Layers (เดี๋ยวเราค่อยมาเติม Code ในไฟล์จริง)
	userRepo := repository.NewPostgresRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 4. Start Server
	app := fiber.New()

	http.NewUserHandler(app, userUsecase)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("E-Wallet API is Running! 🚀")
	})

	// ใช้ Port จาก Config
	log.Fatal(app.Listen(":" + cfg.Server.Port))
}
