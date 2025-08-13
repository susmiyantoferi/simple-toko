package repository

import (
	"context"
	"simple-toko/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) (*entity.Product, error)
	Delete(ctx context.Context, id uint) error
	FindById(ctx context.Context, id uint) (*entity.Product, error)
	FindAll(ctx context.Context) ([]*entity.Product, error)
	AddStock(ctx context.Context, id uint, stock int) (*entity.Product, error)
	ReduceStock(ctx context.Context, id uint, stock int) (*entity.Product, error)
	UpdateImage(ctx context.Context, id uint, img string) (*entity.Product, error)
}
