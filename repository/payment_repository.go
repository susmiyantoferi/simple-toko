package repository

import (
	"context"
	"simple-toko/entity"
)

type PaymentRepository interface {
	UploadPayment(ctx context.Context, pym *entity.Payment) (*entity.Payment, error)
	UpdateStatus(ctx context.Context, pym *entity.Payment) (*entity.Payment, error)
	FindById(ctx context.Context, id uint) (*entity.Payment, error)
	FindByOrderId(ctx context.Context, orderId uint) (*entity.Payment, error)
	FindAll(ctx context.Context, page, pageSize int) ([]*entity.Payment,int64, error)
	Delete(ctx context.Context, id uint) error
	//Confirm(ctx context.Context, orderId uint, pym *entity.Payment) (*entity.Payment, error)
	//UpdatePayment(ctx context.Context, pym *entity.Payment) (*entity.Payment, error)
}
