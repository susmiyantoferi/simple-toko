package entity

import (
	"time"

	"gorm.io/gorm"
)

type OrderProduct struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	OrderID   uint           `gorm:"notnull;uniqueIndex:idx_order_product"`
	Order     Order          `gorm:"foreignKey:OrderID;references:ID;OnDelete:RESTRICT;"`
	ProductID uint           `gorm:"notnull;uniqueIndex:idx_order_product"`
	Product   Product        `gorm:"foreignKey:ProductID;references:ID;OnDelete:RESTRICT;"`
	Qty       int            `gorm:"notnull"`
	UnitPrice float64        `gorm:"notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
