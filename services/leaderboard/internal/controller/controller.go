package leaderboard

import (
	"trailbox/services/leaderboard/internal/model"
	"trailbox/services/leaderboard/internal/repository"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller { return &Controller{repo: r} }

func (c *Controller) Upsert(userID string, score int) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	row := &model.Leaderboard{UserID: uid, Score: score}
	return c.repo.Upsert(row)
}

func (c *Controller) GetTop(limit int) ([]model.Leaderboard, error) {
	return c.repo.GetTop(limit)
}

func (c *Controller) GetByUser(userID string) (*model.Leaderboard, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return c.repo.GetByUser(uid)
}
