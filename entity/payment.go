package entity

import "time"

type Payment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	OrderID   uint      `gorm:"notnull"`
	Order     Order     `gorm:"foreignKey:OrderID;references:ID"`
	Image     string    `gorm:"notnull"`
	CreatedAt time.Time `gorm:"notnull"`
	UpdatedAt time.Time `gorm:"notnull"`
}
