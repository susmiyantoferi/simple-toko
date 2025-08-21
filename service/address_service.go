package service

import (
	"context"
	web "simple-toko/web/address"
	pg "simple-toko/web"
)

type AddressService interface {
	Create(ctx context.Context, req *web.AddressCreateRequest) (*web.AddressResponse, error)
	Update(ctx context.Context, req *web.AddressUpdateRequest) (*web.AddressResponse, error)
	Delete(ctx context.Context, id uint) error
	FindByUserId(ctx context.Context, userId uint) ([]*web.AddressResponse, error)
	FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error)
}