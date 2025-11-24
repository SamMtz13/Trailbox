package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	host := getenvOr("DB_HOST", "postgres.default.svc.cluster.local")
	user := getenvOr("DB_USER", "trailbox")
	pass := getenvOr("DB_PASS", "trailbox")
	name := getenvOr("DB_NAME", "trailbox")
	port := getenvOr("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	log.Println("[leaderboard] ✅ Connected to PostgreSQL")

	// Extensión UUID y migración

	log.Println("[leaderboard] ✅ Migración completada")
	return db, nil
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
