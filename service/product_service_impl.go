package service

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	web "simple-toko/web/product"

	"github.com/go-playground/validator/v10"
)

type productServiceImpl struct {
	ProductRepo   repository.ProductRepository
	InventoryRepo repository.InventoryRepository
	Validate      *validator.Validate
}

func NewProductServiceImpl(productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, validate *validator.Validate) *productServiceImpl {
	return &productServiceImpl{
		ProductRepo:   productRepo,
		InventoryRepo: inventoryRepo,
		Validate:      validate,
	}
}

func (p *productServiceImpl) Create(ctx context.Context, req *web.ProductCreateRequest) (*web.ProductResponse, error) {
	if err := p.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	if _, err := p.InventoryRepo.FindById(ctx, req.InventoryID); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find inventory: %w", err)
	}

	product := entity.Product{
		InventoryID: req.InventoryID,
		Name:        req.Name,
		Price:       req.Price,
		Stock:       req.Stock,
		Description: req.Description,
	}
	result, err := p.ProductRepo.Create(ctx, &product)
	if err != nil {
		return nil, fmt.Errorf("product service: create: %w", err)
	}

	response := helper.ToProductResponse(result)
	return response, nil
}

func (p *productServiceImpl) Update(ctx context.Context, req *web.ProductUpdateRequest) (*web.ProductResponse, error) {
	if err := p.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	prod, err := p.ProductRepo.FindById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find product: %w", err)
	}

	if req.InventoryID != nil {
		if _, err := p.InventoryRepo.FindById(ctx, *req.InventoryID); err != nil {
			if errors.Is(err, repository.ErrorIdNotFound) {
				return nil, ErrorIdNotFound
			}
			return nil, fmt.Errorf("product service: find inventory: %w", err)
		}

		prod.InventoryID = *req.InventoryID
	}

	if req.Name != nil {
		prod.Name = *req.Name
	}

	if req.Price != nil {
		prod.Price = *req.Price
	}

	if req.Description != nil {
		prod.Description = *req.Description
	}

	result, err := p.ProductRepo.Update(ctx, prod)
	if err != nil {
		return nil, fmt.Errorf("product service: update: %w", err)
	}

	response := helper.ToProductResponse(result)
	return response, nil
}

func (p *productServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := p.ProductRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return ErrorIdNotFound
		}
		return fmt.Errorf("product service: delete: %w", err)
	}

	return nil
}

func (p *productServiceImpl) FindById(ctx context.Context, id uint) (*web.ProductResponse, error) {
	result, err := p.ProductRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find id: %w", err)
	}
	response := helper.ToProductResponse(result)
	return response, nil
}

func (p *productServiceImpl) FindAll(ctx context.Context) ([]*web.ProductResponse, error) {
	result, err := p.ProductRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("product service: find all: %w", err)
	}

	var responses []*web.ProductResponse
	for _, v := range result {
		response := web.ProductResponse{
			ID:          v.ID,
			InventoryID: v.InventoryID,
			Inventory: web.InventInfo{
				Location: v.Inventory.Location,
			},
			Name:        v.Name,
			Price:       v.Price,
			Stock:       v.Stock,
			Description: v.Description,
			Image:       v.Image,
			CreatedAt:   v.CreatedAt,
			UpdatedAt:   v.UpdatedAt,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (p *productServiceImpl) AddStock(ctx context.Context, req *web.ProductStockUpdateRequest) (*web.ProductResponse, error) {

	if err := p.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	_, err := p.ProductRepo.FindById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find id add stock: %w", err)
	}

	result, err := p.ProductRepo.AddStock(ctx, req.ID, req.Stock)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: add stock: %w", err)
	}

	response := helper.ToProductResponse(result)
	return response, nil
}

func (p *productServiceImpl) ReduceStock(ctx context.Context, req *web.ProductStockUpdateRequest) (*web.ProductResponse, error) {
	if err := p.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	_, err := p.ProductRepo.FindById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find id reduce stock: %w", err)
	}

	result, err := p.ProductRepo.ReduceStock(ctx, req.ID, req.Stock)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		if errors.Is(err, repository.ErrNotEnoughStock) {
			return nil, ErrNotEnoughStock
		}
		return nil, fmt.Errorf("product service: reduce stock: %w", err)
	}

	response := helper.ToProductResponse(result)
	return response, nil
}

func (p *productServiceImpl) UpdateImage(ctx context.Context, id uint, img string) (*web.ProductResponse, error) {
	result, err := p.ProductRepo.UpdateImage(ctx, id, img)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: update img: %w", err)
	}

	response := helper.ToProductResponse(result)
	return response, nil
}
