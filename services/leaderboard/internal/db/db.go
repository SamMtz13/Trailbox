package db

import (
	"fmt"
	"log"
	"os"

	"trailbox/services/leaderboard/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	log.Println("[leaderboard] ✅ Connected to PostgreSQL")

	// Extensión UUID y migración
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Printf("[leaderboard] ⚠️ could not ensure uuid-ossp: %v", err)
	}
	if err := db.AutoMigrate(&model.Leaderboard{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Println("[leaderboard] ✅ Migración completada")
	return db, nil
}
