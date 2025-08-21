package service

import (
	"context"
	"fmt"
	"math"
	"simple-toko/entity"
	"simple-toko/helper"
	"simple-toko/repository"
	"simple-toko/web"
)

type reportServiceImpl struct {
	ReportRepository repository.ReportRepository
}

func NewReportServiceImpl(reportRepository repository.ReportRepository) *reportServiceImpl {
	return &reportServiceImpl{
		ReportRepository: reportRepository,
	}
}

func (r *reportServiceImpl) MonthlySales(ctx context.Context, page, pageSize int) (*web.PaginatedResponse, error) {
	result, totalItems, err := r.ReportRepository.MonthlySales(ctx, page, pageSize)
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

	totalPage := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	paginateResp := helper.ToPaginatedResponse(int64(page), totalPage, totalItems, responses)
	
	return paginateResp, nil
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
