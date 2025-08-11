package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type addressRepositoryImpl struct {
	Db *gorm.DB
}

func NewAddressRepositoryImpl(db *gorm.DB) *addressRepositoryImpl {
	return &addressRepositoryImpl{
		Db: db,
	}
}

func (a *addressRepositoryImpl) Create(ctx context.Context, adrs *entity.Address) (*entity.Address, error) {
	if err := a.Db.WithContext(ctx).Create(adrs).Error; err != nil {
		return nil, fmt.Errorf("address repo: create: %w", err)
	}

	if err := a.Db.WithContext(ctx).Preload("User").First(adrs, adrs.ID).Error; err != nil{
		return nil, fmt.Errorf("address repo: preload user: %w", err)
	}
	return adrs, nil
}

func (a *addressRepositoryImpl) Update(ctx context.Context, adrs *entity.Address) (*entity.Address, error) {
	data := entity.Address{
		Addresses: adrs.Addresses,
	}

	result := a.Db.WithContext(ctx).First(adrs, adrs.ID).Updates(data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("address repo: update: %w", result.Error)
	}

	if err := a.Db.WithContext(ctx).Preload("User").First(adrs, adrs.ID).Error; err != nil{
		return nil, fmt.Errorf("address repo: preload user: %w", err)
	}

	return adrs, nil
}

func (a *addressRepositoryImpl) Delete(ctx context.Context, id uint) error {
	data := entity.Address{}
	result := a.Db.WithContext(ctx).Delete(&data, id)
	if result.Error != nil {
		return fmt.Errorf("address repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil
}

func (a *addressRepositoryImpl) FindByUserId(ctx context.Context, usrId uint) ([]*entity.Address, error) {
	var data []*entity.Address

	result := a.Db.WithContext(ctx).Preload("User").Where("user_id = ?", usrId).Find(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("address repo: find by user id: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, ErrorIdNotFound
	}

	return data, nil
}

func (a *addressRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Address, error) {
	var data []*entity.Address

	result := a.Db.WithContext(ctx).Preload("User").Find(&data)
	if result.Error != nil {
		return nil, fmt.Errorf("address repo: find all: %w", result.Error)
	}

	return data, nil
}

func (a *addressRepositoryImpl) FindByIdAndUserId(ctx context.Context, id, usrId uint) (*entity.Address, error) {
	data := entity.Address{}

	result := a.Db.WithContext(ctx).Preload("User").Where("id = ? AND user_id = ?", id, usrId).First(&data)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("address repo: find by id: %w", result.Error)
	}

	return &data, nil
}
