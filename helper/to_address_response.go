package helper

import (
	"simple-toko/entity"
	web "simple-toko/web/address"
)

func ToAddressResponse(adrs *entity.Address) *web.AddressResponse{
	return &web.AddressResponse{
		User: web.UserInfo{
			Name: adrs.User.Name,
			Email: adrs.User.Email,
		},
		ID: adrs.ID,
		Addresses: adrs.Addresses,
		CreatedAt: adrs.CreatedAt,
		UpdatedAt: adrs.UpdatedAt,
	}
}