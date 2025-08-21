package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	pg "simple-toko/web"
	web "simple-toko/web/address"

	"github.com/go-playground/validator/v10"
)

type addressServiceImpl struct {
	AddressRepository repository.AddressRepository
	UserRepository    repository.UserRepository
	Validate          *validator.Validate
}

func NewAddressServiceImpl(addressRepository repository.AddressRepository, userRepository repository.UserRepository, validate *validator.Validate) *addressServiceImpl {
	return &addressServiceImpl{
		AddressRepository: addressRepository,
		UserRepository:    userRepository,
		Validate:          validate,
	}
}

func (a *addressServiceImpl) Create(ctx context.Context, req *web.AddressCreateRequest) (*web.AddressResponse, error) {
	if err := a.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	usrId := req.UserID

	_, err := a.UserRepository.FindById(ctx, usrId)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("address service: find user: %w", err)
	}

	adrs := entity.Address{
		UserID:    usrId,
		Addresses: req.Addresses,
	}

	result, err := a.AddressRepository.Create(ctx, &adrs)
	if err != nil {
		return nil, fmt.Errorf("address service: create: %w", err)
	}
	response := helper.ToAddressResponse(result)

	return response, nil
}

func (a *addressServiceImpl) Update(ctx context.Context, req *web.AddressUpdateRequest) (*web.AddressResponse, error) {
	if err := a.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	adrs, err := a.AddressRepository.FindByIdAndUserId(ctx, req.ID, req.UserID)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("address service: update: %w", err)
	}

	adrs.Addresses = req.Addresses

	result, err := a.AddressRepository.Update(ctx, adrs)
	if err != nil {
		return nil, fmt.Errorf("address service: update: %w", err)
	}

	response := helper.ToAddressResponse(result)
	return response, nil
}

func (a *addressServiceImpl) Delete(ctx context.Context, id uint) error {

	if err := a.AddressRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return ErrorIdNotFound
		}
		return fmt.Errorf("address service: delete: %w", err)
	}

	return nil
}

func (a *addressServiceImpl) FindByUserId(ctx context.Context, userId uint) ([]*web.AddressResponse, error) {
	result, err := a.AddressRepository.FindByUserId(ctx, userId)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("address service: find user id: %w", err)
	}

	var responses []*web.AddressResponse
	for _, v := range result {
		response := web.AddressResponse{
			User: web.UserInfo{
				Name:  v.User.Name,
				Email: v.User.Email,
			},
			ID:        v.ID,
			Addresses: v.Addresses,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (a *addressServiceImpl) FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error) {
	result, totalItems, err := a.AddressRepository.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("address service: find all: %w", err)
	}

	var responses []*web.AddressResponse
	for _, v := range result {
		response := web.AddressResponse{
			User: web.UserInfo{
				Name:  v.User.Name,
				Email: v.User.Email,
			},
			ID:        v.ID,
			Addresses: v.Addresses,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)

	return paginateResp, nil
}
