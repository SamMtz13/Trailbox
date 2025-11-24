package db

import (
	"fmt"
	"log"
	"os"

	"trailbox/services/workouts/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Connect() (*gorm.DB, error) {
	host := getenvOr("DB_HOST", "postgres.final-project.svc.cluster.local")
	user := getenvOr("DB_USER", "trailbox")
	pass := getenvOr("DB_PASS", "trailbox")
	name := getenvOr("DB_NAME", "trailbox")
	port := getenvOr("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=America/Mexico_City",
		host, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	log.Println("[workouts] ✅ Connected to PostgreSQL")

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Printf("[workouts] ⚠️ could not ensure uuid-ossp: %v", err)
	}

	if err := db.AutoMigrate(&model.Workout{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	log.Println("[workouts] ✅ Migración completada")
	return db, nil
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
