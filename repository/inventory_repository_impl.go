package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type inventoryRepositoryImpl struct {
	Db *gorm.DB
}

func NewInventoryRepositoryImpl(db *gorm.DB) *inventoryRepositoryImpl {
	return &inventoryRepositoryImpl{
		Db: db,
	}
}

func (i *inventoryRepositoryImpl) Create(ctx context.Context, inv *entity.Inventory) (*entity.Inventory, error) {
	if err := i.Db.WithContext(ctx).Create(inv).Error; err != nil {
		return nil, fmt.Errorf("inventory repo: create: %w", err)
	}
	return inv, nil
}

func (i *inventoryRepositoryImpl) Update(ctx context.Context, invId uint, inv *entity.Inventory) (*entity.Inventory, error) {
	var dataInv entity.Inventory
	if err := i.Db.WithContext(ctx).First(&dataInv, invId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("inventory repo: update stock: %w", err)
	}

	update := map[string]interface{}{
		"location" : inv.Location,
	}

	result := i.Db.WithContext(ctx).Model(&dataInv).Updates(update)
	if result.Error != nil {
		return nil, fmt.Errorf("inventory repo: update stock: %w", result.Error)
	}

	return &dataInv, nil
}

func (i *inventoryRepositoryImpl) Delete(ctx context.Context, invId uint) error {
	dataInv := entity.Inventory{}
	result := i.Db.WithContext(ctx).Delete(&dataInv, invId)
	if result.Error != nil {
		return fmt.Errorf("inventory repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil
}

func (i *inventoryRepositoryImpl) FindById(ctx context.Context, invId uint) (*entity.Inventory, error) {
	dataInv := entity.Inventory{}

	result := i.Db.WithContext(ctx).First(&dataInv, invId)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}

		return nil, fmt.Errorf("inventory repo: find id: %w", result.Error)
	}

	return &dataInv, nil
}

func (i *inventoryRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Inventory, error) {
	var dataInv []*entity.Inventory
	result := i.Db.WithContext(ctx).Find(&dataInv)
	if result.Error != nil {
		return nil, fmt.Errorf("inventory repo: find all: %w", result.Error)
	}

	return dataInv, nil
}
