package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trailbox/services/notifications/internal/model"
)

type Repository struct {
	db *gorm.DB
}

// Conexión a PostgreSQL
func Connect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || name == "" || port == "" {
		return nil, fmt.Errorf("missing database environment variables")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	log.Println("[notifications] ✅ Connected to PostgreSQL")
	return db, nil
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Inserta una nueva notificación
func (r *Repository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

// Lista las notificaciones de un usuario
func (r *Repository) ListByUser(userID string) ([]*model.Notification, error) {
	var notifications []*model.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}
