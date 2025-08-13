package web

import (
	"time"
)

type InventInfo struct {
	Location string `json:"location"`
}

type ProductResponse struct {
	ID          uint       `json:"id"`
	InventoryID uint       `json:"inventory_id"`
	Inventory   InventInfo `json:"inventory"`
	Name        string     `json:"name"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	Description string     `json:"description"`
	Image       string     `json:"image"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
