package helper

import (
	"simple-toko/entity"
	web "simple-toko/web/user"
)

func ToUserResponse(user *entity.User) *web.UserResponse {
	return &web.UserResponse{
		Name: user.Name,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}