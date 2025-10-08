package db

import (
	"context"

	"trailbox/services/routes/internal/model"
	"trailbox/services/routes/internal/repository"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) repository.Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRoute(ctx context.Context, route *model.Route) error {
	return r.db.WithContext(ctx).Create(route).Error
}

func (r *Repository) GetRoute(id string) (*model.Route, error) {
	var route model.Route
	if err := r.db.First(&route, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &route, nil
}

func (r *Repository) ListRoutes() ([]model.Route, error) {
	var routes []model.Route
	if err := r.db.Find(&routes).Error; err != nil {
		return nil, err
	}
	return routes, nil
}
