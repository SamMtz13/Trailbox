package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Age       int       `gorm:"not null"`
	Email     string    `gorm:"type:varchar(200);uniqueIndex;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
