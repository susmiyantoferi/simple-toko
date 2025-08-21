package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"simple-toko/helper"
	"simple-toko/service"
	web "simple-toko/web/payment"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type paymentHandlerImpl struct {
	PaymentService service.PaymentService
}

func NewPaymentHandlerImpl(paymentService service.PaymentService) *paymentHandlerImpl {
	return &paymentHandlerImpl{
		PaymentService: paymentService,
	}
}

var Path = "uploads/payment/"

func (pay *paymentHandlerImpl) UploadPayment(ctx *gin.Context) {
	req := web.PaymentCreateRequest{}

	if err := ctx.ShouldBind(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type file", err.Error())
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	if err := ctx.SaveUploadedFile(file, Path+fileName); err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "failed upload file", err.Error())
		return
	}

	req.Image = fileName

	result, err := pay.PaymentService.UploadPayment(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			os.Remove(Path + fileName)
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrOrderNotFound):
			os.Remove(Path + fileName)
			helper.ToResponseJson(ctx, http.StatusNotFound, "order not found", err.Error())
			return
		default:
			os.Remove(Path + fileName)
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "failed save file", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (pay *paymentHandlerImpl) UpdateStatus(ctx *gin.Context) {
	req := web.PaymentUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("orderId")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input tpye id", err.Error())
		return
	}

	req.OrderID = uint(orderId)

	result, err := pay.PaymentService.UpdateStatus(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "order not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "updated", result)
}

func (pay *paymentHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("id")
	payId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input tpye id", err.Error())
		return
	}

	result, err := pay.PaymentService.FindById(ctx, uint(payId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "payment not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (pay *paymentHandlerImpl) FindByOrderId(ctx *gin.Context) {
	id := ctx.Param("orderId")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input tpye id", err.Error())
		return
	}

	result, err := pay.PaymentService.FindByOrderId(ctx, uint(orderId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "order not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (pay *paymentHandlerImpl) FindAll(ctx *gin.Context) {

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

	result, err := pay.PaymentService.FindAll(ctx, page, pageSize)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (pay *paymentHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	payId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input tpye id", err.Error())
		return
	}

	if err := pay.PaymentService.Delete(ctx, uint(payId)); err != nil {
		switch {
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "payment not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "deleted", nil)
}

func (pay *paymentHandlerImpl) PreviewImage(ctx *gin.Context) {
	id := ctx.Param("orderId")
	orderId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input tpye id", err.Error())
		return
	}

	result, err := pay.PaymentService.FindByOrderId(ctx, uint(orderId))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			helper.ToResponseJson(ctx, http.StatusNotFound, "order not found", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	if result.Image == "" {
		helper.ToResponseJson(ctx, http.StatusNotFound, "image not found", nil)
		return
	}

	file := Path + result.Image
	download := ctx.DefaultQuery("download", "false")

	if download == "true" {
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", result.Image))
	}

	ctx.File(file)
}
