package entity

import "time"

type Order struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	UserID         uint           `gorm:"notnull"`
	User           User           `gorm:"foreignKey:UserID;references:ID"`
	AddressID      uint           `gorm:"notnull"`
	Address        Address        `gorm:"foreignKey:AddressID;references:ID"`
	OrderProducts  []OrderProduct `gorm:"foreignKey:OrderID"`
	StatusOrder    string         `gorm:"notnull"`
	StatusDelivery string         `gorm:"notnull"`
	CreatedAt      time.Time      `gorm:"notnull"`
	UpdatedAt      time.Time      `gorm:"notnull"`
}
