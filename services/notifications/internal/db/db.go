package db

import (
	"fmt"
	"log"
	"os"

	"trailbox/services/notifications/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

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

	if err := db.AutoMigrate(&model.Notification{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("[notifications] ✅ Connected to PostgreSQL")
	return db, nil
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

func (r *Repository) ListByUser(userID string) ([]*model.Notification, error) {
	var notifications []*model.Notification
	err := r.db.Where("user_id = ?", userID).Find(&notifications).Error
	return notifications, err
}
