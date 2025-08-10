package entity

import (
	"time"

	"gorm.io/gorm"
)

type Inventory struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Location  string         `gorm:"size:100;notnull"`
	Stock     int            `gorm:"notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
