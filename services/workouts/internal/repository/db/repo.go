package db

import (
	"trailbox/services/workouts/internal/model"
	"trailbox/services/workouts/internal/repository"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

// ✅ Constructor que te faltaba
func New(conn *gorm.DB) repository.Repository {
	return &DBRepository{db: conn}
}

// Crear workout
func (r *DBRepository) Create(w *model.Workout) error {
	return r.db.Create(w).Error
}

// Obtener por ID
func (r *DBRepository) GetByID(id uuid.UUID) (*model.Workout, error) {
	var workout model.Workout
	if err := r.db.First(&workout, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &workout, nil
}

// Listar todos
func (r *DBRepository) List() ([]*model.Workout, error) {
	var workouts []*model.Workout
	if err := r.db.Find(&workouts).Error; err != nil {
		return nil, err
	}
	return workouts, nil
}
