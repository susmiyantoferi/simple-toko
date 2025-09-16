package helper

import (
	"simple-toko/entity"
	web "simple-toko/web/payment"
)

func ToPaymentResponse(pay *entity.Payment) *web.PaymentResponse {
	return &web.PaymentResponse{
		OrderID: pay.OrderID,
		Order: web.OrderInfo{
			AmountPay:      pay.Order.AmountPay,
			StatusOrder:    pay.Order.StatusOrder,
			StatusDelivery: pay.Order.StatusDelivery,
		},
		Image:     pay.Image,
		Status:    pay.Status,
		CreatedAt: pay.CreatedAt,
		UpdatedAt: pay.UpdatedAt,
	}
}
