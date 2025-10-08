package model

import (
	"time"

	"github.com/google/uuid"
)

// Route representa una ruta creada por un usuario.
type Route struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Path      string    `gorm:"type:varchar(255);not null"`
	Duration  int       `gorm:"not null"` // duración en minutos
	Distance  int       `gorm:"not null"` // distancia en metros o km (según uses)
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
