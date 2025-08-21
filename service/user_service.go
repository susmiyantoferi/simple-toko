package service

import (
	"context"
	token "simple-toko/web"
	web "simple-toko/web/user"
)

type UserService interface {
	Create(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error)
	Update(ctx context.Context, userId uint, req *web.UserUpdateRequest) (*web.UserResponse, error)
	Delete(ctx context.Context, userId uint) error
	FindById(ctx context.Context, userId uint) (*web.UserResponse, error)
	FindByEmail(ctx context.Context, email string) (*web.UserResponse, error)
	FindAll(ctx context.Context, page, pageSize int) (*token.PaginatedResponse, error)
	CreateAdmin(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error)
	Login(ctx context.Context, req *web.UserLoginRequest) (*token.TokenResponse, error)
	RefreshToken(ctx context.Context, req *web.UserRefreshTokenRequest) (*token.TokenResponse, error)
}
