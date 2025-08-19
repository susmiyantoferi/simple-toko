package repository

import (
	"simple-toko/entity"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type reportRepositoryImpl struct {
	Db *gorm.DB
}

func NewReportRepositoryImpl(db *gorm.DB) *reportRepositoryImpl {
	return &reportRepositoryImpl{
		Db: db,
	}
}

func (r *reportRepositoryImpl) MonthlySales(ctx context.Context) ([]*entity.SalesReport, error) {
	var data []*entity.SalesReport

	err := r.Db.WithContext(ctx).Table("order_products AS op").
		Select("DATE_FORMAT(o.created_at, '%Y-%m') AS month, SUM(op.qty) AS total_qty, SUM(op.qty * op.unit_price) AS total_sales").
		Joins("JOIN orders o ON op.order_id = o.id").Where("o.status_order = ?", Confirmed).
		Group("DATE_FORMAT(o.created_at, '%Y-%m')").Order("month").Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *reportRepositoryImpl) TopProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error) {
	var data []*entity.TopProduct

	err := r.Db.WithContext(ctx).Table("order_products AS op").
		Select("p.id AS product_id, p.name AS product_name, SUM(op.qty) AS total_qty, SUM(op.qty * op.unit_price) AS total_sales").
		Joins("JOIN products p ON op.product_id = p.id").Joins("JOIN orders o ON op.order_id = o.id").
		Where("o.status_order = ?", Confirmed).Group("p.id, p.name").Order("total_qty DESC").Limit(limit).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *reportRepositoryImpl) LessProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error) {
	var data []*entity.TopProduct

	err := r.Db.WithContext(ctx).Table("order_products AS op").
		Select("p.id AS product_id, p.name AS product_name, SUM(op.qty) AS total_qty, SUM(op.qty * op.unit_price) AS total_sales").
		Joins("JOIN products p ON op.product_id = p.id").Joins("JOIN orders o ON op.order_id = o.id").
		Where("o.status_order = ?", Confirmed).Group("p.id, p.name").Order("total_qty ASC").Limit(limit).Scan(&data).Error

	if err != nil {
		return nil, err
	}

	return data, nil
}
