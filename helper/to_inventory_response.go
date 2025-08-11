package helper

import (
	"simple-toko/entity"
	web "simple-toko/web/inventory"
)

func ToInventoryResponse(inv *entity.Inventory) *web.InventoryResponse{
	return &web.InventoryResponse{
		Location: inv.Location,
		CreatedAt: inv.CreatedAt,
		UpdatedAt: inv.UpdatedAt,
	}
}