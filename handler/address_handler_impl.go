package handler

import (
	"errors"
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	t "simple-toko/web"
	web "simple-toko/web/address"
	"strconv"

	"github.com/gin-gonic/gin"
)

type addressHandlerImpl struct {
	AddressService service.AddressService
}

func NewAddressHandlerImpl(addressService service.AddressService) *addressHandlerImpl {
	return &addressHandlerImpl{
		AddressService: addressService,
	}
}

func (a *addressHandlerImpl) Create(ctx *gin.Context) {
	req := web.AddressCreateRequest{}

	claims, _ := ctx.Get("user")
	user := claims.(*t.TokenClaim)

	req.UserID = user.UserID

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := a.AddressService.Create(ctx, &req)
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

func (a *addressHandlerImpl) Update(ctx *gin.Context) {
	req := web.AddressUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
	}

	adrsId := ctx.Param("id")
	addressId, err := strconv.Atoi(adrsId)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	claims, _ := ctx.Get("user")
	user := claims.(*t.TokenClaim)

	req.ID = uint(addressId)
	req.UserID = user.UserID

	result, err := a.AddressService.Update(ctx, &req)
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

func (a *addressHandlerImpl) Delete(ctx *gin.Context) {
	adrsId := ctx.Param("id")
	addressId, err := strconv.Atoi(adrsId)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	if err := a.AddressService.Delete(ctx, uint(addressId)); err != nil {
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

func (a *addressHandlerImpl) FindByUserId(ctx *gin.Context) {
	usrId := ctx.Param("userId")
	userId, err := strconv.Atoi(usrId)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := a.AddressService.FindByUserId(ctx, uint(userId))
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

func (a *addressHandlerImpl) FindAll(ctx *gin.Context) {

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

	result, err := a.AddressService.FindAll(ctx, page, pageSize)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}
	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}
