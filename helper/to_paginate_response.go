package helper

import "simple-toko/web"

func ToPaginatedResponse(currentPage int64, totalPage int, totalItems int64, data interface{}) *web.PaginatedResponse {
	return &web.PaginatedResponse{
		CurrentPage: currentPage,
		TotalPage:   totalPage,
		TotalItems:  totalItems,
		Data:        data,
	}
}
