package repository

import (
	"context"
	"simple-toko/entity"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order *entity.Order) (*entity.Order, error)
	UpdateAddress(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.Order, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.Order, int64, error)
	FindByOrderId(ctx context.Context, orderId uint) ([]*entity.OrderProduct, error)
	ConfirmOrder(ctx context.Context, orderId uint, order *entity.Order) (*entity.Order, error)
	// RemoveOrderItem(ctx context.Context, orderId, productId uint) (*entity.Order, error)
	// UpdateOrderQty(ctx context.Context, orderId, productId uint, qty int) (*entity.Order, error)
	//AddOrderItem(ctx context.Context, orderId uint, item *entity.Order) (*entity.Order, error)

}
