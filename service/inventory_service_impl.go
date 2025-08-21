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
	web "simple-toko/web/inventory"

	"github.com/go-playground/validator/v10"
)

type inventoryServiceImpl struct {
	InventoryRepo repository.InventoryRepository
	Validate      *validator.Validate
}

func NewInventoryServiceImpl(inventoryRepo repository.InventoryRepository, validate *validator.Validate) *inventoryServiceImpl {
	return &inventoryServiceImpl{
		InventoryRepo: inventoryRepo,
		Validate:      validate,
	}
}

var ErrNotEnoughStock = errors.New("not enough stock")

func (i *inventoryServiceImpl) Create(ctx context.Context, req *web.InventoryCreateRequest) (*web.InventoryResponse, error) {
	if err := i.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	inv := entity.Inventory{
		Location: req.Location,
	}

	result, err := i.InventoryRepo.Create(ctx, &inv)
	if err != nil {
		return nil, fmt.Errorf("user service: create: %w", err)
	}

	response := helper.ToInventoryResponse(result)

	return response, nil
}

func (i *inventoryServiceImpl) Update(ctx context.Context, invId uint, req *web.InventoryCreateRequest) (*web.InventoryResponse, error) {

	if err := i.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	update := entity.Inventory{
		Location: req.Location,
	}

	result, err := i.InventoryRepo.Update(ctx, invId, &update)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	response := helper.ToInventoryResponse(result)

	return response, nil

}

func (i *inventoryServiceImpl) Delete(ctx context.Context, invId uint) error {
	if err := i.InventoryRepo.Delete(ctx, invId); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return ErrorIdNotFound
		}
		return fmt.Errorf("user service: delete: %w", err)
	}

	return nil
}

func (i *inventoryServiceImpl) FindById(ctx context.Context, invId uint) (*web.InventoryResponse, error) {

	result, err := i.InventoryRepo.FindById(ctx, invId)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}

		return nil, fmt.Errorf("user service: find id: %w", err)
	}

	response := helper.ToInventoryResponse(result)
	return response, nil

}

func (i *inventoryServiceImpl) FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error) {
	result, totalItems, err := i.InventoryRepo.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("user service: find all: %w", err)
	}

	var responses []*web.InventoryResponse
	for _, v := range result {
		response := web.InventoryResponse{
			Location:  v.Location,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}

		responses = append(responses, &response)
	}

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)

	return paginateResp, nil
}
