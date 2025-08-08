package entity

import "time"

type Inventory struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Location  string    `gorm:"notnull"`
	Stock     int       `gorm:"notnull"`
	CreatedAt time.Time `gorm:"notnull"`
	UpdatedAt time.Time `gorm:"notnull"`
}
