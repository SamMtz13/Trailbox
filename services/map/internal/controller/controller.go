package mapctrl

import (
	"trailbox/services/map/internal/model"
	"trailbox/services/map/internal/repository"

	"github.com/google/uuid"
)

type Controller struct {
	repo repository.Repository
}

func NewController(r repository.Repository) *Controller {
	return &Controller{repo: r}
}

func (c *Controller) SetRouteMap(routeID string, geoJSON string) error {
	rid, err := uuid.Parse(routeID)
	if err != nil {
		return err
	}
	return c.repo.SetRouteMap(rid, geoJSON)
}

func (c *Controller) GetRouteMap(routeID string) (*model.Map, error) {
	rid, err := uuid.Parse(routeID)
	if err != nil {
		return nil, err
	}
	return c.repo.GetByRouteID(rid)
}

func (c *Controller) ListMaps() ([]model.Map, error) {
	return c.repo.List()
}
