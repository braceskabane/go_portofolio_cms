package database

import (
	"log"
	"portfolio-cms/internal/config"
	"portfolio-cms/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) *gorm.DB {
	logLevel := logger.Silent
	if cfg.App.Env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// Connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get underlying sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	DB = db
	log.Println("✅ Database connected successfully")
	return db
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Skill{},
		&model.Experience{},
		&model.Education{},
		&model.Contact{},
		&model.Profile{},
	)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("✅ Database migrated successfully")
}
