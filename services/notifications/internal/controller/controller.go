package controller

import (
	"trailbox/services/notifications/internal/model"
	"trailbox/services/notifications/internal/repository/db"
)

// Controller contiene la lógica de negocio.
type Controller struct {
	repo *db.Repository
}

func NewController(r *db.Repository) *Controller {
	return &Controller{repo: r}
}

// Crea una nueva notificación
func (c *Controller) Create(userID, message string) (*model.Notification, error) {
	n := &model.Notification{
		UserID:  userID,
		Message: message,
	}
	if err := c.repo.Create(n); err != nil {
		return nil, err
	}
	return n, nil
}

// Lista todas las notificaciones de un usuario
func (c *Controller) ListByUser(userID string) ([]*model.Notification, error) {
	return c.repo.ListByUser(userID)
}
