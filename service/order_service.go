package service

import (
	"context"
	web "simple-toko/web/order"
	pg "simple-toko/web"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req *web.OrderCreateRequest) (*web.OrderResponse, error)
	UpdateAddress(ctx context.Context, req *web.OrderUpdateRequest) (*web.OrderResponse, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*web.OrderResponse, error)
	FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error)
	// FindByOrderId(ctx context.Context, id uint) ([]*web.OrderResponse, error)
	ConfirmOrder(ctx context.Context, req *web.OrderUpdateStatusRequest) (*web.OrderResponse, error)
}
