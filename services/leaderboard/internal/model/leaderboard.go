package model

import (
	"time"

	"github.com/google/uuid"
)

type Leaderboard struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"` // FK -> users.id (lógica FK si la quieres después)
	Score     int       `gorm:"not null"`
	Position  int       `gorm:"not null;default:0"` // puedes calcularla al consultar
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
