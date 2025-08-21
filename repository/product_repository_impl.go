package repository

import (
	"context"
	"errors"
	"fmt"
	"simple-toko/entity"

	"gorm.io/gorm"
)

type productRepositoryImpl struct {
	Db *gorm.DB
}

func NewProductRepositoryImpl(db *gorm.DB) *productRepositoryImpl {
	return &productRepositoryImpl{
		Db: db,
	}
}

func (p *productRepositoryImpl) Create(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	if err := p.Db.WithContext(ctx).Create(product).Error; err != nil {
		return nil, fmt.Errorf("product repo: create: %w", err)
	}

	if err := p.Db.WithContext(ctx).Preload("Inventory").First(product, product.ID).Error; err != nil {
		return nil, fmt.Errorf("product repo: preload create: %w", err)
	}

	return product, nil
}

func (p *productRepositoryImpl) Update(ctx context.Context, product *entity.Product) (*entity.Product, error) {
	data := entity.Product{
		InventoryID: product.InventoryID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
	}

	if err := p.Db.WithContext(ctx).First(product, product.ID).Updates(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product repo: update: %w", err)
	}

	if err := p.Db.WithContext(ctx).Preload("Inventory").First(product, product.ID).Error; err != nil {
		return nil, fmt.Errorf("product repo: preload update: %w", err)
	}

	return product, nil
}

func (p *productRepositoryImpl) Delete(ctx context.Context, id uint) error {
	data := entity.Product{}

	result := p.Db.WithContext(ctx).Delete(&data, id)
	if result.Error != nil {
		return fmt.Errorf("product repo: delete: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return ErrorIdNotFound
	}

	return nil
}

func (p *productRepositoryImpl) FindById(ctx context.Context, id uint) (*entity.Product, error) {
	product := entity.Product{}

	if err := p.Db.WithContext(ctx).Preload("Inventory").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product repo: find by id: %w", err)
	}

	return &product, nil
}

func (p *productRepositoryImpl) FindAll(ctx context.Context, page, pageSize int, search string) ([]*entity.Product, int64, error) {
	var product []*entity.Product
	var totalItems int64

	query := p.Db.WithContext(ctx).Model(&entity.Product{})

	if search != ""{
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := query.Preload("Inventory").Limit(pageSize).
		Offset(offset).Find(&product).Error; err != nil {
		return nil, 0, err
	}

	return product, totalItems, nil
}

func (p *productRepositoryImpl) AddStock(ctx context.Context, id uint, stock int) (*entity.Product, error) {

	if stock <= 0 {
		return nil, ErrorValidation
	}

	data := map[string]interface{}{
		"stock": gorm.Expr("stock + ?", stock),
	}

	result := p.Db.WithContext(ctx).Model(&entity.Product{}).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return nil, fmt.Errorf("product repo: add stock: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, ErrorIdNotFound
	}

	var newProd entity.Product
	if err := p.Db.WithContext(ctx).Preload("Inventory").First(&newProd, id).Error; err != nil {
		return nil, fmt.Errorf("product repo: preload add stock: %w", err)
	}

	return &newProd, nil
}

func (p *productRepositoryImpl) ReduceStock(ctx context.Context, id uint, stock int) (*entity.Product, error) {
	if stock <= 0 {
		return nil, ErrorValidation
	}

	data := map[string]interface{}{
		"stock": gorm.Expr("stock - ?", stock),
	}

	var currentStock int
	if err := p.Db.WithContext(ctx).Model(&entity.Product{}).Select("stock").Where("id = ?", id).Scan(&currentStock).Error; err != nil {
		return nil, fmt.Errorf("product repo: get stock: %w", err)
	}

	if currentStock < stock {
		return nil, ErrNotEnoughStock
	}

	result := p.Db.WithContext(ctx).Model(&entity.Product{}).Where("id = ? AND stock >= ?", id, stock).Updates(data)
	if result.Error != nil {
		return nil, fmt.Errorf("product repo: reduce stock: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, ErrorIdNotFound
	}

	var newProd entity.Product
	if err := p.Db.WithContext(ctx).Preload("Inventory").First(&newProd, id).Error; err != nil {
		return nil, fmt.Errorf("product repo: preload reduce stock: %w", err)
	}

	return &newProd, nil
}

func (p *productRepositoryImpl) UpdateImage(ctx context.Context, id uint, img string) (*entity.Product, error) {
	product := entity.Product{}

	if err := p.Db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrorIdNotFound
		}
		return nil, fmt.Errorf("product repo: find id image: %w", err)
	}

	result := p.Db.WithContext(ctx).Model(&product).Update("image", img)
	if result.Error != nil {
		return nil, fmt.Errorf("product repo: update image: %w", result.Error)
	}

	var data entity.Product
	if err := p.Db.WithContext(ctx).Preload("Inventory").First(&data, id).Error; err != nil {
		return nil, fmt.Errorf("product repo: preload update img: %w", err)
	}

	return &data, nil
}
