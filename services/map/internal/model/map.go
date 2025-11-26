package model

import (
	"time"

	"github.com/google/uuid"
)

type Map struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:id"`
	RouteID   uuid.UUID `gorm:"type:uuid;not null;column:route_id"` // FK -> routes.id
	GeoJSON   string    `gorm:"type:jsonb;not null;column:geojson"` // 👈 OJO: geojson
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at"`
}
