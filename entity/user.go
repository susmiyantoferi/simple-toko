package entity

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"notnull"`
	Email     string    `gorm:"notnull"`
	Password  string    `gorm:"notnull"`
	Role      string    `gorm:"notnull"`
	CreatedAt time.Time `gorm:"notnull"`
	UpdatedAt time.Time `gorm:"notnull"`
}
