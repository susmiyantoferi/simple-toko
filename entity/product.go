package entity

import "time"

type Product struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	InventoryID   uint           `gorm:"notnull"`
	Inventory     Inventory      `gorm:"foreignKey:InventoryID;references:ID"`
	Name          string         `gorm:"notnull"`
	Price         float64        `gorm:"notnull"`
	Description   string         `gorm:"notnull"`
	Image         string         `gorm:"notnull"`
	OrderProducts []OrderProduct `gorm:"foreignKey:ProductID"`
	CreatedAt     time.Time      `gorm:"notnull"`
	UpdatedAt     time.Time      `gorm:"notnull"`
}
