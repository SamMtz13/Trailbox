package controller

import (
	"context"
	"sort"
	"time"
)

type Entry struct {
	UserID        string    `json:"user_id"`
	DistanceTotal float64   `json:"distance_km_total"`
	Workouts      int       `json:"workouts_count"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Repository interface {
	Upsert(ctx context.Context, e Entry) (Entry, error)
	Top(ctx context.Context, limit int) ([]Entry, error)
}

type Controller struct{ repo Repository }

func NewController(r Repository) *Controller { return &Controller{repo: r} }

func (c *Controller) Upsert(ctx context.Context, userID string, distance float64, workouts int) (Entry, error) {
	e := Entry{
		UserID:        userID,
		DistanceTotal: distance,
		Workouts:      workouts,
		UpdatedAt:     time.Now().UTC(),
	}
	return c.repo.Upsert(ctx, e)
}

func (c *Controller) Top(ctx context.Context, limit int) ([]Entry, error) {
	all, err := c.repo.Top(ctx, limit)
	if err != nil {
		return nil, err
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].DistanceTotal == all[j].DistanceTotal {
			return all[i].Workouts > all[j].Workouts
		}
		return all[i].DistanceTotal > all[j].DistanceTotal
	})
	if limit > 0 && len(all) > limit {
		return all[:limit], nil
	}
	return all, nil
}
