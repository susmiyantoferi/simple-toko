package entity

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	OrderID   uint           `gorm:"notnull"`
	Order     Order          `gorm:"foreignKey:OrderID;references:ID"`
	Image     string         `gorm:"size:255;notnull"`
	Status    string         `gorm:"type:enum('waiting','confirmed','canceled');default:'waiting';notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
