package repository

import (
	"context"
	"simple-toko/entity"
)

type ReportRepository interface {
	MonthlySales(ctx context.Context) ([]*entity.SalesReport, error)
	TopProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
	LessProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
}
