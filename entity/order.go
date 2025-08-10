package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID             uint           `gorm:"primaryKey;autoIncrement"`
	UserID         uint           `gorm:"notnull"`
	User           User           `gorm:"foreignKey:UserID;references:ID"`
	AddressID      uint           `gorm:"notnull"`
	Address        Address        `gorm:"foreignKey:AddressID;references:ID"`
	OrderProducts  []OrderProduct `gorm:"foreignKey:OrderID"`
	StatusOrder    string         `gorm:"type:enum('waiting','confirmed','canceled');default:'waiting';notnull"`
	StatusDelivery string         `gorm:"type:enum('waiting','on process','delivered','canceled');default:'waiting';notnull"`
	CreatedAt      time.Time      `gorm:"notnull"`
	UpdatedAt      time.Time      `gorm:"notnull"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
