package handler

import (
	"errors"
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	web "simple-toko/web/inventory"
	"strconv"

	"github.com/gin-gonic/gin"
)

type inventoryHandlerImpl struct {
	InventoryService service.InventoryService
}

func NewInventoryHandlerImpl(inventoryService service.InventoryService) *inventoryHandlerImpl {
	return &inventoryHandlerImpl{
		InventoryService: inventoryService,
	}
}

func (i *inventoryHandlerImpl) Create(ctx *gin.Context) {
	req := web.InventoryCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := i.InventoryService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (i *inventoryHandlerImpl) Update(ctx *gin.Context) {
	req := web.InventoryCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("invId")
	invId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := i.InventoryService.Update(ctx, uint(invId), &req)
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

func (i *inventoryHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("invId")
	invId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	if err := i.InventoryService.Delete(ctx, uint(invId)); err != nil {
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

func (i *inventoryHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("invId")
	invId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := i.InventoryService.FindById(ctx, uint(invId))
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

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)

}

func (i *inventoryHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := i.InventoryService.FindAll(ctx)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}
