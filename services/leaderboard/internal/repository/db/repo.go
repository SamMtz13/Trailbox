package db

import (
	"errors"

	"trailbox/services/leaderboard/internal/model"
	"trailbox/services/leaderboard/internal/repository"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

func New(conn *gorm.DB) repository.Repository {
	return &DBRepository{db: conn}
}

func (r *DBRepository) Upsert(lb *model.Leaderboard) error {
	var existing model.Leaderboard
	err := r.db.Where("user_id = ?", lb.UserID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(lb).Error
	}
	if err != nil {
		return err
	}
	existing.Score = lb.Score
	// Position la puedes recalcular luego
	return r.db.Save(&existing).Error
}

func (r *DBRepository) GetTop(limit int) ([]model.Leaderboard, error) {
	var rows []model.Leaderboard
	if err := r.db.Order("score DESC, created_at ASC").Limit(limit).Find(&rows).Error; err != nil {
		return nil, err
	}
	// Asigna posición 1..n en memoria
	for i := range rows {
		rows[i].Position = i + 1
	}
	return rows, nil
}

func (r *DBRepository) GetByUser(userID uuid.UUID) (*model.Leaderboard, error) {
	var row model.Leaderboard
	if err := r.db.Where("user_id = ?", userID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
