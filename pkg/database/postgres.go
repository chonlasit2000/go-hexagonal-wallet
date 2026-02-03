package database

import (
	"fmt"
	"log"
	"time"

	"github.com/chonlasit2000/e-wallet-hexagonal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDatabase(cfg *config.Config) *gorm.DB {
	// DSN ของ Postgres หน้าตาประมาณนี้
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.Port,
		cfg.Database.SSLMode,
		cfg.Database.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // ให้ปริ้น SQL ออกมาดูตอน Dev
	})

	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	// จูนตามความเหมาะสมของ Server
	sqlDB.SetMaxIdleConns(10)  // Connection ที่เปิดรอไว้
	sqlDB.SetMaxOpenConns(100) // Connection สูงสุดที่รับได้
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db

}
