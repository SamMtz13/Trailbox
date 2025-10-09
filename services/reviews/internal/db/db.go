package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establece conexión con PostgreSQL
func Connect() (*gorm.DB, error) {
	host := getenvOr("DB_HOST", "postgres")
	port := getenvOr("DB_PORT", "5432")
	user := getenvOr("DB_USER", "trailbox")
	pass := getenvOr("DB_PASS", "trailbox")
	name := getenvOr("DB_NAME", "trailbox")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=America/Mexico_City",
		host, port, user, pass, name,
	)

	log.Printf("[reviews][db] 📦 Conectando a PostgreSQL en %s:%s...", host, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("[reviews][db] ❌ error al conectar a PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("[reviews][db] ❌ error obteniendo handle SQL: %w", err)
	}

	// Configuración básica de conexión
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0)

	log.Println("[reviews][db] ✅ Conexión establecida correctamente")

	return db, nil
}

// getenvOr obtiene variable de entorno o valor por defecto
func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
