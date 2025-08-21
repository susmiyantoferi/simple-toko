package handler

import (
	"errors"
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	t "simple-toko/web"
	web "simple-toko/web/order"
	"strconv"

	"github.com/gin-gonic/gin"
)

type orderHandlerImpl struct {
	OrderService service.OrderService
}

func NewOrderHandlerImpl(orderService service.OrderService) *orderHandlerImpl {
	return &orderHandlerImpl{
		OrderService: orderService,
	}
}

func (o *orderHandlerImpl) CreateOrder(ctx *gin.Context) {
	req := web.OrderCreateRequest{}

	claims, _ := ctx.Get("user")
	user := claims.(*t.TokenClaim)

	req.UserID = user.UserID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := o.OrderService.CreateOrder(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAddress):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "cannot use address", err.Error())
			return
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrAddressNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "address not found", err.Error())
			return
		case errors.Is(err, service.ErrProductNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "product not found", err.Error())
			return
		case errors.Is(err, service.ErrNotEnoughStock):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "stock not enough", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (o *orderHandlerImpl) UpdateAddress(ctx *gin.Context) {
	req := web.OrderUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	req.ID = uint(orderId)

	result, err := o.OrderService.UpdateAddress(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAddress):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "cannot use address", err.Error())
			return
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "order not found", err.Error())
			return
		case errors.Is(err, service.ErrAddressNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "address not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "updated", result)
}

func (o *orderHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	if err := o.OrderService.Delete(ctx, uint(orderId)); err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "orders not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "deleted", nil)
}

func (o *orderHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := o.OrderService.FindById(ctx, uint(orderId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "orders not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (o *orderHandlerImpl) FindAll(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 5
	}

	result, err := o.OrderService.FindAll(ctx, page, pageSize)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (o *orderHandlerImpl) ConfirmOrder(ctx *gin.Context) {
	req := web.OrderUpdateStatusRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("id")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	req.ID = uint(orderId)

	result, err := o.OrderService.ConfirmOrder(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "orders not found", err.Error())
			return
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}
