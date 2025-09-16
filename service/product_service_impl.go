package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	"simple-toko/utils"
	pg "simple-toko/web"
	web "simple-toko/web/product"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type productServiceImpl struct {
	ProductRepo   repository.ProductRepository
	InventoryRepo repository.InventoryRepository
	Validate      *validator.Validate
	Redis         *redis.Client
}

func NewProductServiceImpl(productRepo repository.ProductRepository, inventoryRepo repository.InventoryRepository, validate *validator.Validate, redis *redis.Client) *productServiceImpl {
	return &productServiceImpl{
		ProductRepo:   productRepo,
		InventoryRepo: inventoryRepo,
		Validate:      validate,
		Redis:         redis,
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

	utils.InvalidateCached(ctx, p.Redis, result.ID)

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

	utils.InvalidateCached(ctx, p.Redis, result.ID)

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

	utils.InvalidateCached(ctx, p.Redis, id)

	return nil
}

func (p *productServiceImpl) FindById(ctx context.Context, id uint) (*web.ProductResponse, error) {
	cacheKey := fmt.Sprintf("products:%d", id)

	cached, err := p.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var response web.ProductResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}
	} else if err != redis.Nil {
		fmt.Printf("Redis error: %v\n", err)
	}

	result, err := p.ProductRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product service: find id: %w", err)
	}
	response := helper.ToProductResponse(result)

	jsonData, _ := json.Marshal(response)
	err = p.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		fmt.Printf("Redis set error: %v\n", err)
	}

	return response, nil
}

func (p *productServiceImpl) FindAll(ctx context.Context, page, pageSize int, search string) (*pg.PaginatedResponse, error) {
	cacheKey := fmt.Sprintf("products:page=%d:size=%d:search=%s", page, pageSize, search)

	cached, err := p.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var paginateResp pg.PaginatedResponse
		if err := json.Unmarshal([]byte(cached), &paginateResp); err == nil {
			return &paginateResp, nil
		}
	} else if err != redis.Nil {
		fmt.Printf("Redis error: %v\n", err)
	}

	result, totalItems, err := p.ProductRepo.FindAll(ctx, page, pageSize, search)
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

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)

	jsonData, _ := json.Marshal(paginateResp)
	err = p.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err()
	if err != nil {
		fmt.Printf("Redis set error: %v\n", err)
	}

	return paginateResp, nil
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

	utils.InvalidateCached(ctx, p.Redis, result.ID)

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


	utils.InvalidateCached(ctx, p.Redis, result.ID)

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

	utils.InvalidateCached(ctx, p.Redis, result.ID)

	response := helper.ToProductResponse(result)
	return response, nil
}
