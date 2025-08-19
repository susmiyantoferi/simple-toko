package service

import (
	"context"
	"fmt"
	"simple-toko/entity"
	"simple-toko/repository"
)

type reportServiceImpl struct {
	ReportRepository repository.ReportRepository
}

func NewReportServiceImpl(reportRepository repository.ReportRepository) *reportServiceImpl {
	return &reportServiceImpl{
		ReportRepository: reportRepository,
	}
}

func (r *reportServiceImpl) MonthlySales(ctx context.Context) ([]*entity.SalesReport, error) {
	result, err := r.ReportRepository.MonthlySales(ctx)
	if err != nil {
		return nil, fmt.Errorf("report service: monthly sales: %w", err)
	}

	var responses []*entity.SalesReport
	for _, v := range result {
		response := entity.SalesReport{
			Month:      v.Month,
			TotalQty:   v.TotalQty,
			TotalSales: v.TotalSales,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (r *reportServiceImpl) TopProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error) {
	result, err := r.ReportRepository.TopProductSales(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("report service: top product: %w", err)
	}

	var responses []*entity.TopProduct
	for _, v := range result {
		response := entity.TopProduct{
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			TotalQty:    v.TotalQty,
			TotalSales:  v.TotalSales,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}

func (r *reportServiceImpl) LessProductSales(ctx context.Context, limit int) ([]*entity.TopProduct, error) {
	result, err := r.ReportRepository.LessProductSales(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("report service: less product: %w", err)
	}

	var responses []*entity.TopProduct
	for _, v := range result {
		response := entity.TopProduct{
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			TotalQty:    v.TotalQty,
			TotalSales:  v.TotalSales,
		}
		responses = append(responses, &response)
	}

	return responses, nil
}
