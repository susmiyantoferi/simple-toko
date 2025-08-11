package repository

import (
	"context"
	"simple-toko/entity"
)

type InventoryRepository interface {
	Create(ctx context.Context, inv *entity.Inventory) (*entity.Inventory, error)
	Update(ctx context.Context, invId uint, inv *entity.Inventory) (*entity.Inventory, error)
	Delete(ctx context.Context, invId uint) error
	FindById(ctx context.Context, invId uint) (*entity.Inventory, error)
	FindAll(ctx context.Context) ([]*entity.Inventory, error)
}
