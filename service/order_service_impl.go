package service

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	adrs "simple-toko/web/address"
	web "simple-toko/web/order"

	"github.com/go-playground/validator/v10"
)

type orderServiceImpl struct {
	OrderRepository repository.OrderRepository
	Validate        *validator.Validate
}

func NewOrderServiceImpl(orderRepository repository.OrderRepository, validate *validator.Validate) *orderServiceImpl {
	return &orderServiceImpl{
		OrderRepository: orderRepository,
		Validate:        validate,
	}
}

var (
	ErrEmptyItems      = errors.New("order has no items")
	ErrProductNotFound = errors.New("product not found")
	ErrOrderNotFound   = errors.New("order not found")
	ErrAddressNotFound = errors.New("address not found")
)

func (o *orderServiceImpl) CreateOrder(ctx context.Context, req *web.OrderCreateRequest) (*web.OrderResponse, error) {
	if err := o.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
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
		if errors.Is(err, repository.ErrAddressNotFound) {
			return nil, ErrAddressNotFound
		}

		if errors.Is(err, repository.ErrProductNotFound) {
			return nil, ErrProductNotFound
		}

		if errors.Is(err, repository.ErrNotEnoughStock) {
			return nil, ErrNotEnoughStock
		}

		return nil, fmt.Errorf("order service: create order: %w", err)
	}

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

	order.AddressID = req.AddressID

	result, err := o.OrderRepository.UpdateAddress(ctx, order)
	if err != nil {
		if errors.Is(err, repository.ErrAddressNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("order service: update address: %w", err)
	}

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
	return nil
}

func (o *orderServiceImpl) FindById(ctx context.Context, id uint) (*web.OrderResponse, error) {
	result, err := o.OrderRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: find by id order: %w", err)
	}

	response := helper.ToOrderResponse(result)
	return response, nil
}

func (o *orderServiceImpl) FindAll(ctx context.Context) ([]*web.OrderResponse, error) {
	result, err := o.OrderRepository.FindAll(ctx)
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
					Price:       op.Product.Price,
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

	return responses, nil
}

// func (o orderServiceImpl) FindByOrderId(ctx context.Context, id uint) ([]*web.OrderResponse, error) {

// }

func (o *orderServiceImpl) ConfirmOrder(ctx context.Context, req *web.OrderUpdateStatusRequest) (*web.OrderResponse, error) {
	if err := o.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	result, err := o.OrderRepository.ConfirmOrder(ctx, req.ID, req.StatusOrder, req.StatusDelivery)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("order service: find id order confirm: %w", err)
	}	

	response := helper.ToOrderResponse(result)
	return response, nil
}
