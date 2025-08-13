package repository

import (
	"context"
	"simple-toko/entity"
)

type AddressRepository interface {
	Create(ctx context.Context, addrss *entity.Address) (*entity.Address, error)
	Update(ctx context.Context, addrss *entity.Address) (*entity.Address, error)
	Delete(ctx context.Context, id uint) error
	FindByUserId(ctx context.Context, userId uint) ([]*entity.Address, error)
	FindAll(ctx context.Context) ([]*entity.Address, error)
	FindByIdAndUserId(ctx context.Context, id, usrId uint) (*entity.Address, error)
}
