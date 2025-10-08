package repository

import (
	"trailbox/services/leaderboard/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	Upsert(lb *model.Leaderboard) error
	GetTop(limit int) ([]model.Leaderboard, error)
	GetByUser(userID uuid.UUID) (*model.Leaderboard, error)
}
