package helper

import (
	"simple-toko/entity"
	web "simple-toko/web/product"
)

func ToProductResponse(product *entity.Product) *web.ProductResponse {
	return &web.ProductResponse{
		ID:          product.ID,
		InventoryID: product.InventoryID,
		Inventory: web.InventInfo{
			Location: product.Inventory.Location,
		},
		Name:        product.Name,
		Price:       product.Price,
		Stock:       product.Stock,
		Description: product.Description,
		Image:       product.Image,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}
