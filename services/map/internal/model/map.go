package model

import (
	"time"

	"github.com/google/uuid"
)

type Map struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RouteID   uuid.UUID `gorm:"type:uuid;not null"`  // FK -> routes.id
	GeoJSON   string    `gorm:"type:jsonb;not null"` // datos geográficos
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
