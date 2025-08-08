package entity

import "time"

type OrderProduct struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	OrderID   uint      `gorm:"notnull"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID"`
	ProductID uint      `gorm:"notnull"`
	Product   Product   `gorm:"foreignKey:ProductID;references:ID"`
	Qty       int       `gorm:"notnull"`
	UnitPrice float64   `gorm:"notnull"`
	CreatedAt time.Time `gorm:"notnull"`
	UpdatedAt time.Time `gorm:"notnull"`
}
