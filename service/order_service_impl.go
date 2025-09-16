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
	adrs "simple-toko/web/address"
	web "simple-toko/web/order"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)

type orderServiceImpl struct {
	OrderRepository  repository.OrderRepository
	AddressRepostory repository.AddressRepository
	Validate         *validator.Validate
	Redis            *redis.Client
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository, addressRepostory repository.AddressRepository, validate *validator.Validate, redis *redis.Client) *orderServiceImpl {
	return &orderServiceImpl{
		OrderRepository:  orderRepository,
		AddressRepostory: addressRepostory,
		Validate:         validate,
		Redis:            redis,
	}
}

var (
	ErrEmptyItems      = errors.New("order has no items")
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrAddressNotFound = errors.New("address not found")
	ErrInvalidAddress  = errors.New("invalid input address")
)

func (o *orderServiceImpl) CreateOrder(ctx context.Context, req *web.OrderCreateRequest) (*web.OrderResponse, error) {
	if err := o.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	_, err := o.AddressRepostory.FindByIdAndUserId(ctx, req.AddressID, req.UserID)
	if err != nil {
		return nil, ErrAddressNotFound
	}

	order := entity.Order{
		UserID:        req.UserID,
		AddressID:     req.AddressID,
		OrderProducts: make([]entity.OrderProduct, len(req.OrderProducts)),
	}

	for i, v := range req.OrderProducts {
		order.OrderProducts[i] = entity.OrderProduct{
			ProductID: v.ProductID,
			Qty:       v.Qty,
		}
	}

	result, err := o.OrderRepository.CreateOrder(ctx, &order)
	if err != nil {

		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}

		if errors.Is(err, repository.ErrNotEnoughStock) {
			return nil, ErrNotEnoughStock
		}

		return nil, fmt.Errorf("order service: create order: %w", err)
	}

	utils.InvalidateCached(ctx, o.Redis, result.ID)

	response := helper.ToOrderResponse(result)

	return response, nil
}

func (o *orderServiceImpl) UpdateAddress(ctx context.Context, req *web.OrderUpdateRequest) (*web.OrderResponse, error) {
	if err := o.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	order, err := o.OrderRepository.FindById(ctx, req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: find order update address: %w", err)
	}

	if _, err := o.AddressRepostory.FindByIdAndUserId(ctx, req.AddressID, order.UserID); err != nil {
		return nil, ErrAddressNotFound
	}

	order.AddressID = req.AddressID

	result, err := o.OrderRepository.UpdateAddress(ctx, order)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: update address: %w", err)
	}

	utils.InvalidateCached(ctx, o.Redis, result.ID)

	response := helper.ToOrderResponse(result)
	return response, nil
}

func (o *orderServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := o.OrderRepository.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return ErrorIdNotFound
		}
		return fmt.Errorf("order service: delete order: %w", err)
	}

	utils.InvalidateCached(ctx, o.Redis, id)
	return nil
}

func (o *orderServiceImpl) FindById(ctx context.Context, id uint) (*web.OrderResponse, error) {
	cacheKey := fmt.Sprintf("orders:%d", id)

	cached, err := o.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var response web.OrderResponse
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			return &response, nil
		}

	} else if err != redis.Nil {
		fmt.Printf("Redis error: %v\n", err)
	}

	result, err := o.OrderRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: find by id order: %w", err)
	}

	response := helper.ToOrderResponse(result)

	jsonData, _ := json.Marshal(response)
	if err := o.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err(); err != nil {
		fmt.Printf("Redis set error: %v\n", err)
	}
	return response, nil
}

func (o *orderServiceImpl) FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error) {
	cacheKey := fmt.Sprintf("orders:page=%d:size=%d", page, pageSize)
	cached, err := o.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var paginateResp pg.PaginatedResponse
		if err := json.Unmarshal([]byte(cached), &paginateResp); err == nil {
			return &paginateResp, nil
		}
	} else if err != redis.Nil {
		fmt.Printf("Redis error: %v\n", err)
	}

	result, totalItems, err := o.OrderRepository.FindAll(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("order service: find all order: %w", err)
	}

	var responses []*web.OrderResponse
	for _, v := range result {
		orderProducts := []web.OrderProductInfo{}
		for _, op := range v.OrderProducts {
			orderProducts = append(orderProducts, web.OrderProductInfo{
				ProductID: op.ProductID,
				Product: web.ProductInfo{
					Name:        op.Product.Name,
					Description: op.Product.Description,
					Image:       op.Product.Image,
				},
				Qty:       op.Qty,
				UnitPrice: op.UnitPrice,
			})
		}
		response := web.OrderResponse{
			ID:        v.ID,
			AmountPay: v.AmountPay,
			User: adrs.UserInfo{
				Name:  v.User.Name,
				Email: v.User.Email,
			},
			AddressID: v.AddressID,
			Address: web.AddressInfo{
				Addresses: v.Address.Addresses,
			},
			OrderProducts:  orderProducts,
			StatusOrder:    v.StatusOrder,
			StatusDelivery: v.StatusDelivery,
			CreatedAt:      v.CreatedAt,
			UpdatedAt:      v.UpdatedAt,
		}

		responses = append(responses, &response)
	}

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)

	jsonData, _ := json.Marshal(paginateResp)
	if err := o.Redis.Set(ctx, cacheKey, jsonData, 5*time.Minute).Err(); err != nil {
		fmt.Printf("Redis set error: %v\n", err)
	}

	return paginateResp, nil
}

// func (o orderServiceImpl) FindByOrderId(ctx context.Context, id uint) ([]*web.OrderResponse, error) {

// }

func (o *orderServiceImpl) ConfirmOrder(ctx context.Context, req *web.OrderUpdateStatusRequest) (*web.OrderResponse, error) {
	if err := o.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	order := entity.Order{}

	if req.StatusOrder != nil {
		order.StatusOrder = *req.StatusOrder
	}

	if req.StatusDelivery != nil {
		order.StatusDelivery = *req.StatusDelivery
	}

	result, err := o.OrderRepository.ConfirmOrder(ctx, req.ID, &order)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: find id order confirm: %w", err)
	}

	utils.InvalidateCached(ctx, o.Redis, result.ID)

	response := helper.ToOrderResponse(result)
	return response, nil
}
