package db

import (
	"trailbox/services/map/internal/model"
	"trailbox/services/map/internal/repository"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type DBRepository struct {
	db *gorm.DB
}

func New(conn *gorm.DB) repository.Repository {
	return &DBRepository{db: conn}
}

func (r *DBRepository) SetRouteMap(routeID uuid.UUID, geoJSON string) error {
	var existing model.Map
	err := r.db.Where("route_id = ?", routeID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		newMap := model.Map{RouteID: routeID, GeoJSON: geoJSON}
		return r.db.Create(&newMap).Error
	}
	if err != nil {
		return err
	}
	existing.GeoJSON = geoJSON
	return r.db.Save(&existing).Error
}

func (r *DBRepository) GetByRouteID(routeID uuid.UUID) (*model.Map, error) {
	var m model.Map
	if err := r.db.Where("route_id = ?", routeID).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *DBRepository) List() ([]model.Map, error) {
	var maps []model.Map
	if err := r.db.Find(&maps).Error; err != nil {
		return nil, err
	}
	return maps, nil
}
