package db

import (
	"context"

	"trailbox/services/users/internal/model"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUser(id string) (*model.User, error) {
	var u model.User
	if err := r.db.First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) ListUsers() ([]model.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *Repository) CreateUser(ctx context.Context, u *model.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}
