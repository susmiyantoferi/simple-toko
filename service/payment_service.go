package service

import (
	"context"
	web "simple-toko/web/payment"
)

type PaymentService interface {
	UploadPayment(ctx context.Context, req *web.PaymentCreateRequest) (*web.PaymentResponse, error)
	UpdateStatus(ctx context.Context, req *web.PaymentUpdateRequest) (*web.PaymentResponse, error)
	FindById(ctx context.Context, id uint) (*web.PaymentResponse, error)
	FindByOrderId(ctx context.Context, orerId uint) (*web.PaymentResponse, error)
	FindAll(ctx context.Context) ([]*web.PaymentResponse, error)
	Delete(ctx context.Context, id uint) error
	//UpdatePayment(ctx context.Context, req *web.PaymentUpdateRequest) (*web.PaymentResponse, error)
}
