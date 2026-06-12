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

	// Auto-migrate models
	if err := AutoMigrate(db); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("✅ Database migrated successfully")

	DB = db
	log.Println("✅ Database connected successfully")
	return db
}

// AutoMigrate runs GORM auto-migration for all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Auth
		&model.User{},

		// Profile
		&model.Profile{},

		// Tech Stack — category harus sebelum item (FK)
		&model.StackCategory{},
		&model.StackItem{},

		// Project Category harus sebelum Project
		&model.ProjectCategory{},

		// Experience Category harus sebelum Experience
		&model.ExperienceCategory{},

		// Skill (tidak ada FK constraint ke tabel lain)
		&model.Skill{},

		// Project & Experience (relasi many2many akan dibuat otomatis)
		&model.Project{},
		&model.Experience{},

		// Education
		&model.Education{},

		// Contact
		&model.Contact{},

		// Running Activity
		&model.RunningActivity{},

		// Asset (polymorphic, aman setelah tabel owner)
		&model.Asset{},
	)
}
