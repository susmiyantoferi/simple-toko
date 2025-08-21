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
	web "simple-toko/web/payment"

	"github.com/go-playground/validator/v10"
)

type paymentServiceImpl struct {
	PaymentRepo repository.PaymentRepository
	OrderRepo repository.OrderRepository
	Validate    *validator.Validate
}

func NewPaymentServiceImpl(paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, validate *validator.Validate) *paymentServiceImpl {
	return &paymentServiceImpl{
		PaymentRepo: paymentRepo,
		OrderRepo: orderRepo,
		Validate:    validate,
	}
}

func (pay *paymentServiceImpl) UploadPayment(ctx context.Context, req *web.PaymentCreateRequest) (*web.PaymentResponse, error) {
	if err := pay.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	orderId := req.OrderID

	_, err := pay.OrderRepo.FindById(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment service: find order, upload pay: %w", err)
	}

	data := entity.Payment{
		OrderID: req.OrderID,
		Image:   req.Image,
		Status: repository.Waiting,
	}

	result, err := pay.PaymentRepo.UploadPayment(ctx, &data)
	if err != nil {
		return nil, fmt.Errorf("payment service: upload payment: %w", err)
	}

	response := helper.ToPaymentResponse(result)
	return response, nil
}

func (pay *paymentServiceImpl) UpdateStatus(ctx context.Context, req *web.PaymentUpdateRequest) (*web.PaymentResponse, error) {
	if err := pay.Validate.Struct(req); err != nil {
		return nil, ErrorValidation
	}

	orderId := req.OrderID

	order, err := pay.PaymentRepo.FindByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment service: find order, update status: %w", err)
	}

	if req.Status != nil {
		order.Status = *req.Status
	}

	result, err := pay.PaymentRepo.UpdateStatus(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("payment service: update status: %w", err)
	}

	response := helper.ToPaymentResponse(result)
	return response, nil
}

func (pay *paymentServiceImpl) FindById(ctx context.Context, id uint) (*web.PaymentResponse, error) {
	result, err := pay.PaymentRepo.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("payment service: find payment: %w", err)
	}

	response := helper.ToPaymentResponse(result)
	return response, nil
}

func (pay *paymentServiceImpl) FindByOrderId(ctx context.Context, orderId uint) (*web.PaymentResponse, error) {
	result, err := pay.PaymentRepo.FindByOrderId(ctx, orderId)
	if err != nil {
		if errors.Is(err, repository.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment service: find order: %w", err)
	}

	response := helper.ToPaymentResponse(result)
	return response, nil
}

func (pay *paymentServiceImpl) FindAll(ctx context.Context, page, pageSize int) (*pg.PaginatedResponse, error) {
	result, totalItems, err := pay.PaymentRepo.FindAll(ctx, page, pageSize)

	if err != nil {
		return nil, fmt.Errorf("payment service: find all: %w", err)
	}

	var responses []*web.PaymentResponse
	for _, v := range result {
		response := web.PaymentResponse{
			OrderID: v.OrderID,
			Order: web.OrderInfo{
				AmountPay:      v.Order.AmountPay,
				StatusOrder:    v.Order.StatusOrder,
				StatusDelivery: v.Order.StatusDelivery,
			},
			Image:     v.Image,
			Status:    v.Status,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		}

		responses = append(responses, &response)
	}

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)

	return paginateResp, nil
}

func (pay *paymentServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := pay.PaymentRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrorIdNotFound) {
			return ErrorIdNotFound
		}
		return fmt.Errorf("payment service: delete: %w", err)
	}

	return nil
}
