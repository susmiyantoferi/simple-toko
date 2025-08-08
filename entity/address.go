package entity

import "time"

type Address struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"notnull"`
	User      User      `gorm:"foreignKey:UserID;references:ID"`
	Addresses string    `gorm:"notnull"`
	CreatedAt time.Time `gorm:"notnull"`
	UpdatedAt time.Time `gorm:"notnull"`
}
