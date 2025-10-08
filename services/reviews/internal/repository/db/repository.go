package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"trailbox/services/reviews/internal/model"
)

type Repository struct {
	db *gorm.DB
}

// Conecta con PostgreSQL usando variables de entorno
func Connect() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" || user == "" || name == "" || port == "" {
		return nil, fmt.Errorf("missing one or more DB environment variables")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect DB: %w", err)
	}

	if err := db.AutoMigrate(&model.Review{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	log.Println("[reviews] ✅ Connected to PostgreSQL")
	return db, nil
}

// Crea un nuevo repositorio
func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Crea una nueva review
func (r *Repository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

// Lista las reseñas por ruta
func (r *Repository) ListByRoute(routeID string) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.Where("route_id = ?", routeID).
		Order("created_at DESC").
		Find(&reviews).Error
	return reviews, err
}
