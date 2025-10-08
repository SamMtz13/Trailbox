package controller

import (
	"trailbox/services/reviews/internal/model"
	"trailbox/services/reviews/internal/repository/db"
)

type Controller struct {
	repo *db.Repository
}

func NewController(r *db.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) AddReview(userID, routeID, comment string, rating int) (*model.Review, error) {
	r := &model.Review{
		UserID:  userID,
		RouteID: routeID,
		Comment: comment,
		Rating:  rating,
	}
	if err := c.repo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Controller) ListReviews(routeID string) ([]*model.Review, error) {
	return c.repo.ListByRoute(routeID)
}
