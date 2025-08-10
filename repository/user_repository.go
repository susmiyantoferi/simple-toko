package repository

import (
	"context"
	"simple-toko/entity"
)

type UserRepository interface{
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, userId uint, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, userId uint) error
	FindById(ctx context.Context, userId uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
}