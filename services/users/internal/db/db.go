package db

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = Connect()
		if err != nil {
			log.Fatalf("❌ Could not connect to the database: %v", err)
		}
	})
	return db
}

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
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
