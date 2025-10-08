package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to parse JSONB column")
	}
	return json.Unmarshal(bytes, a)
}

type Workout struct {
	ID        uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string      `gorm:"type:varchar(100);not null"`
	Exercises StringArray `gorm:"type:jsonb"`
	Duration  int         `gorm:"not null"`
	Calories  int         `gorm:"not null"` // 👈 nuevo campo
	Date      time.Time   `gorm:"not null"` // 👈 nuevo campo
	UserID    uuid.UUID   `gorm:"type:uuid;not null"`
	RouteID   uuid.UUID   `gorm:"type:uuid;not null"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
}
