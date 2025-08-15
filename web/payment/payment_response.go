package web

import "time"

type OrderInfo struct {
	AmountPay      float64 `json:"amount_pay"`
	StatusOrder    string  `json:"status_order"`
	StatusDelivery string  `json:"status_delivery"`
}

type PaymentResponse struct {
	OrderID   uint      `json:"order_id"`
	Order     OrderInfo `json:"order"`
	Image     string    `json:"image"`
	Status    string    `json:"status_payment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
