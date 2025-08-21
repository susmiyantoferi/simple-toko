package service

import (
	"context"
	web "simple-toko/web/product"
	pg "simple-toko/web"
)

type ProductService interface {
	Create(ctx context.Context, req *web.ProductCreateRequest) (*web.ProductResponse, error)
	Update(ctx context.Context, req *web.ProductUpdateRequest) (*web.ProductResponse, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*web.ProductResponse, error)
	FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error)
	AddStock(ctx context.Context, req *web.ProductStockUpdateRequest) (*web.ProductResponse, error)
	ReduceStock(ctx context.Context, req *web.ProductStockUpdateRequest) (*web.ProductResponse, error)
	UpdateImage(ctx context.Context, id uint, img string) (*web.ProductResponse, error)
}
