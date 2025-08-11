package service

import (
	"errors"
	"fmt"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	"simple-toko/utils"
	web "simple-toko/web/user"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/net/context"
)

type userServiceImpl struct {
	UserRepository repository.UserRepository
	Validate       *validator.Validate
}

func NewUserServiceImpl(userRepository repository.UserRepository, validate *validator.Validate) *userServiceImpl {
	return &userServiceImpl{
		UserRepository: userRepository,
		Validate:       validate,
	}
}

var ErrorIdNotFound = errors.New("id not found")
var ErrorEmailNotFound = errors.New("email not found")
var ErrorEmailExist = errors.New("email already exist")
var ErrorValidation = errors.New("validation failed")

const Customer string = "customer"
const Admin string = "admin"

func (r *userServiceImpl) Create(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := r.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}

	user := entity.User{
		Name:     req.Name,
		Email:    email,
		Password: pass,
		Role:     Customer,
	}

	created, err := r.UserRepository.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, ErrorEmailExist) {
			return nil, ErrorEmailExist
		}

		return nil, fmt.Errorf("user service: create: %w", err)
	}

	response := helper.ToUserResponse(created)

	return response, nil

}

func (r *userServiceImpl) Update(ctx context.Context, userId uint, req *web.UserUpdateRequest) (*web.UserResponse, error) {
	if err := r.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	user := entity.User{}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Password != nil {
		pass, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("hashing password: %w", err)
		}
		user.Password = pass
	}

	result, err := r.UserRepository.Update(ctx, userId, &user)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: update: %w", err)
	}
	response := helper.ToUserResponse(result)

	return response, nil
}

func (r *userServiceImpl) Delete(ctx context.Context, userId uint) error {
	if err := r.UserRepository.Delete(ctx, userId); err != nil {
		if errors.Is(err, ErrorIdNotFound) {
			return ErrorIdNotFound
		}

		return fmt.Errorf("user service: delete: %w", err)
	}

	return nil
}

func (r *userServiceImpl) FindById(ctx context.Context, userId uint) (*web.UserResponse, error) {

	result, err := r.UserRepository.FindById(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: find id: %w", err)
	}

	response := helper.ToUserResponse(result)

	return response, nil
}

func (r *userServiceImpl) FindByEmail(ctx context.Context, email string) (*web.UserResponse, error) {
	result, err := r.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrorEmailNotFound) {
			return nil, ErrorEmailNotFound
		}

		return nil, fmt.Errorf("user service: find email: %w", err)
	}

	response := helper.ToUserResponse(result)

	return response, nil
}

func (r *userServiceImpl) FindAll(ctx context.Context) ([]*web.UserResponse, error) {

	result, err := r.UserRepository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service: find all: %w", err)
	}

	var responses []*web.UserResponse
	for _, v := range result {
		response := web.UserResponse{
			Name:      v.Name,
			Email:     v.Email,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}

		responses = append(responses, &response)
	}

	return responses, nil
}

func (r *userServiceImpl) CreateAdmin(ctx context.Context, req *web.UserCreateRequest) (*web.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := r.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	pass, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}

	user := entity.User{
		Name:     req.Name,
		Email:    email,
		Password: pass,
		Role:     Admin,
	}

	created, err := r.UserRepository.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, ErrorEmailExist) {
			return nil, ErrorEmailExist
		}

		return nil, ErrorEmailExist
	}

	response := helper.ToUserResponse(created)

	return response, nil

}
