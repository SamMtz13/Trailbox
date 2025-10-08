package repository

import (
	"trailbox/services/map/internal/model"

	"github.com/google/uuid"
)

type Repository interface {
	SetRouteMap(routeID uuid.UUID, geoJSON string) error
	GetByRouteID(routeID uuid.UUID) (*model.Map, error)
	List() ([]model.Map, error)
}
