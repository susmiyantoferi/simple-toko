package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type paymentRepositoryImpl struct {
	Db *gorm.DB
}

func NewPaymentRepositoryImpl(db *gorm.DB) *paymentRepositoryImpl {
	return &paymentRepositoryImpl{
		Db: db,
	}
}

func (pay *paymentRepositoryImpl) UploadPayment(ctx context.Context, pym *entity.Payment) (*entity.Payment, error) {

	if err := pay.Db.WithContext(ctx).First(&entity.Order{}, pym.OrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment repo: find order, upload payment: %w", err)
	}

	if err := pay.Db.WithContext(ctx).Model(pym).Create(pym).Error; err != nil {
		return nil, fmt.Errorf("payment repo: upload payment: %w", err)
	}

	if err := pay.Db.WithContext(ctx).Preload("Order").Where("order_id = ?", pym.OrderID).Take(pym).Error; err != nil {
		return nil, fmt.Errorf("payment repo: preload payment: %w", err)
	}

	return pym, nil
}

func (pay *paymentRepositoryImpl) UpdateStatus(ctx context.Context, pym *entity.Payment) (*entity.Payment, error) {
	if err := pay.Db.WithContext(ctx).First(&entity.Order{}, pym.OrderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment repo: find order, update status: %w", err)
	}

	if err := pay.Db.WithContext(ctx).Model(pym).Where("order_id = ?", pym.OrderID).Update("status", pym.Status).Error; err != nil {
		return nil, fmt.Errorf("payment repo: update status: %w", err)
	}

	if err := pay.Db.WithContext(ctx).Preload("Order").Where("order_id = ?", pym.OrderID).Take(pym).Error; err != nil {
		return nil, fmt.Errorf("payment repo: preload payment: %w", err)
	}

	return pym, nil
}

func (pay *paymentRepositoryImpl) FindById(ctx context.Context, id uint) (*entity.Payment, error) {
	if err := pay.Db.WithContext(ctx).First(&entity.Payment{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("payment repo: find id: %w", err)
	}

	var data entity.Payment
	if err := pay.Db.WithContext(ctx).Preload("Order").First(&data, id).Error; err != nil {
		return nil, fmt.Errorf("payment repo: preload payment: %w", err)
	}

	return &data, nil
}

func (pay *paymentRepositoryImpl) FindByOrderId(ctx context.Context, orderId uint) (*entity.Payment, error) {
	if err := pay.Db.WithContext(ctx).First(&entity.Order{}, orderId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, fmt.Errorf("payment repo: find id: %w", err)
	}

	var data entity.Payment
	if err := pay.Db.WithContext(ctx).Preload("Order").Where("order_id = ?", orderId).Take(&data).Error; err != nil {
		return nil, fmt.Errorf("payment repo: preload payment, find order id: %w", err)
	}

	return &data, nil
}

func (pay *paymentRepositoryImpl) FindAll(ctx context.Context, page, pageSize int) ([]*entity.Payment, int64, error) {
	var dataPay []*entity.Payment
	var totalItems int64

	if err := pay.Db.WithContext(ctx).Model(&entity.Payment{}).Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := pay.Db.WithContext(ctx).Limit(pageSize).Offset(offset).
		Preload("Order").Find(&dataPay).Error; err != nil {
		return nil, 0, nil
	}

	return dataPay, totalItems, nil
}

func (pay *paymentRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := pay.Db.WithContext(ctx).Delete(&entity.Payment{}, id)
	if result.Error != nil {
		return fmt.Errorf("payment repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil
}
