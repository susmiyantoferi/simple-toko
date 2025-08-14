package web

import (
	web "simple-toko/web/address"
	"time"
)

type AddressInfo struct {
	Addresses string `json:"addresses"`
}

type ProductInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
}

type OrderProductInfo struct {
	ProductID uint        `json:"product_id"`
	Product   ProductInfo `json:"product"`
	Qty       int         `json:"qty"`
	UnitPrice float64     `json:"unit_price"`
}

type OrderResponse struct {
	ID             uint               `json:"id"`
	AmountPay      float64            `json:"amount_pay"`
	User           web.UserInfo       `json:"user"`
	AddressID      uint               `json:"address_id"`
	Address        AddressInfo        `json:"address"`
	OrderProducts  []OrderProductInfo `json:"order_product"`
	StatusOrder    string             `json:"status_order"`
	StatusDelivery string             `json:"status_delivery"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
