package entity

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	InventoryID   uint           `gorm:"notnull"`
	Inventory     Inventory      `gorm:"foreignKey:InventoryID;references:ID"`
	Name          string         `gorm:"size:100;notnull"`
	Price         float64        `gorm:"notnull"`
	Stock         int            `gorm:"notnull"`
	Description   string         `gorm:"size:255;notnull"`
	Image         string         `gorm:"size:255;default:null"`
	OrderProducts []OrderProduct `gorm:"foreignKey:ProductID"`
	CreatedAt     time.Time      `gorm:"notnull"`
	UpdatedAt     time.Time      `gorm:"notnull"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
