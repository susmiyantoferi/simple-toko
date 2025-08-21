package service

import (
	"context"
	"simple-toko/entity"
	"simple-toko/web"
)

type ReportService interface {
	MonthlySales(ctx context.Context, page, pageSize int) (*web.PaginatedResponse, error)
	TopProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
	LessProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error)
}
