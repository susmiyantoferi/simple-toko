package service

import (
	"context"
	"simple-toko/entity"
)

type ReportService interface {
	MonthlySales(ctx context.Context) ([]*entity.SalesReport, error)
	TopProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
	LessProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
}
