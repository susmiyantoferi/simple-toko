package handler

import (
	"errors"
	"fmt"
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	web "simple-toko/web/product"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type productHandlerImpl struct {
	ProductService service.ProductService
}

func NewProductHandlerImpl(productService service.ProductService) *productHandlerImpl {
	return &productHandlerImpl{
		ProductService: productService,
	}
}

func (p *productHandlerImpl) Create(ctx *gin.Context) {
	req := web.ProductCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := p.ProductService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (p *productHandlerImpl) Update(ctx *gin.Context) {
	req := web.ProductUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	req.ID = uint(productId)

	result, err := p.ProductService.Update(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "updated", result)

}

func (p *productHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	if err := p.ProductService.Delete(ctx, uint(productId)); err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "deleted", nil)

}

func (p *productHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := p.ProductService.FindById(ctx, uint(productId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (p *productHandlerImpl) FindAll(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "5")
	search := ctx.DefaultQuery("search", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err:= strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1{
		pageSize = 5
	}

	result, err := p.ProductService.FindAll(ctx, page, pageSize, search)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (p *productHandlerImpl) AddStock(ctx *gin.Context) {
	req := web.ProductStockUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	req.ID = uint(productId)

	result, err := p.ProductService.AddStock(ctx, &req)
	if err != nil {
		switch {
			case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input stock", nil)
			return
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (p *productHandlerImpl) ReduceStock(ctx *gin.Context) {
	req := web.ProductStockUpdateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	req.ID = uint(productId)

	result, err := p.ProductService.ReduceStock(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input stock", nil)
			return
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		case errors.Is(err, service.ErrNotEnoughStock):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "not enough stock", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (p *productHandlerImpl) UpdateImage(ctx *gin.Context) {
	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type file", err.Error())
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	if err := ctx.SaveUploadedFile(file, "uploads/product/"+fileName); err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "failed upload file", err.Error())
		return
	}

	result, err := p.ProductService.UpdateImage(ctx, uint(productId), fileName)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "failed save file", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (p *productHandlerImpl) PreviewImage(ctx *gin.Context) {
	id := ctx.Param("productId")
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := p.ProductService.FindById(ctx, uint(productId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}

	if result.Image == "" {
		helper.ToResponseJson(ctx, http.StatusNotFound, "image not found", nil)
		return
	}

	file := "uploads/product/" + result.Image
	download := ctx.DefaultQuery("download", "false")

	if download == "true" {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", result.Image))
	}

	ctx.File(file)
}
