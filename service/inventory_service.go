package service

import (
	"context"
	web "simple-toko/web/inventory"
	pg "simple-toko/web"
)

type InventoryService interface{
	Create(ctx context.Context, req *web.InventoryCreateRequest) (*web.InventoryResponse, error)
	Update(ctx context.Context, invId uint, req *web.InventoryCreateRequest) (*web.InventoryResponse, error)
	Delete(ctx context.Context, invId uint) error
	FindById(ctx context.Context, invId uint) (*web.InventoryResponse, error)
	FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error)
}