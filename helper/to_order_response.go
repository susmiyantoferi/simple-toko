package helper

import (
	"simple-toko/entity"
	adr "simple-toko/web/address"
	web "simple-toko/web/order"
)

func ToOrderResponse(o *entity.Order) *web.OrderResponse {
	orderProducts := make([]web.OrderProductInfo, 0, len(o.OrderProducts))

	for _, v := range o.OrderProducts {
		orderProducts = append(orderProducts, web.OrderProductInfo{
			ProductID: v.ProductID,
			Product: web.ProductInfo{
				Name:        v.Product.Name,
				Price:       v.Product.Price,
				Description: v.Product.Description,
				Image:       v.Product.Image,
			},
			Qty:       v.Qty,
			UnitPrice: v.UnitPrice,
		})
	}

	return &web.OrderResponse{
		ID:        o.ID,
		AmountPay: o.AmountPay,
		User: adr.UserInfo{
			Name:  o.User.Name,
			Email: o.User.Email,
		},
		AddressID: o.AddressID,
		Address: web.AddressInfo{
			Addresses: o.Address.Addresses,
		},
		OrderProducts:  orderProducts,
		StatusOrder:    o.StatusOrder,
		StatusDelivery: o.StatusDelivery,
		CreatedAt:      o.CreatedAt,
		UpdatedAt:      o.UpdatedAt,
	}
}
