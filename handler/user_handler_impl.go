package handler

import (
	"errors"
	"net/http"
	"simple-toko/helper"
	"simple-toko/service"
	web "simple-toko/web/user"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandlerImpl struct {
	UserService service.UserService
}

func NewUserHandlerImpl(userService service.UserService) *userHandlerImpl {
	return &userHandlerImpl{
		UserService: userService,
	}
}

func (h *userHandlerImpl) Create(ctx *gin.Context) {
	req := web.UserCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.UserService.Create(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrorEmailExist):
			helper.ToResponseJson(ctx, http.StatusConflict, "email already exist", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (h *userHandlerImpl) Update(ctx *gin.Context) {
	req := web.UserUpdateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	id := ctx.Param("userId")
	userId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := h.UserService.Update(ctx, uint(userId), &req)
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

func (h *userHandlerImpl) Delete(ctx *gin.Context) {
	id := ctx.Param("userId")
	userId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", err.Error())
		return
	}

	if err := h.UserService.Delete(ctx, uint(userId)); err != nil {
		if errors.Is(err, service.ErrorIdNotFound) {
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", err.Error())
			return
		}
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "deleted", nil)
}

func (h *userHandlerImpl) FindById(ctx *gin.Context) {
	id := ctx.Param("userId")
	userId, err := strconv.Atoi(id)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type id", nil)
		return
	}

	result, err := h.UserService.FindById(ctx, uint(userId))
	if err != nil {
		if errors.Is(err, service.ErrorIdNotFound) {
			helper.ToResponseJson(ctx, http.StatusNotFound, "id not found", err.Error())
			return
		}
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (h *userHandlerImpl) FindByEmail(ctx *gin.Context) {
	email := ctx.Param("email")
	_, err := strconv.Atoi(email)
	if err == nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input type email", nil)
		return
	}

	result, err := h.UserService.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, service.ErrorEmailNotFound) {
			helper.ToResponseJson(ctx, http.StatusNotFound, "email not found", err.Error())
			return
		}
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (h *userHandlerImpl) FindAll(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "5")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1{
		pageSize = 5
	}

	result, err := h.UserService.FindAll(ctx, page, pageSize)
	if err != nil {
		helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (h *userHandlerImpl) CreateAdmin(ctx *gin.Context) {
	req := web.UserCreateRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.UserService.CreateAdmin(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrorEmailExist):
			helper.ToResponseJson(ctx, http.StatusConflict, "email already exist", nil)
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", nil)
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusCreated, "created", result)
}

func (h *userHandlerImpl) Login(ctx *gin.Context) {
	req := web.UserLoginRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.UserService.Login(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrFailedLogin):
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "email or password in correct", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}

func (h *userHandlerImpl) RefreshToken(ctx *gin.Context) {
	req := web.UserRefreshTokenRequest{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid request", err.Error())
		return
	}

	result, err := h.UserService.RefreshToken(ctx, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrorValidation):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid input", err.Error())
			return
		case errors.Is(err, service.ErrInvalidToken):
			helper.ToResponseJson(ctx, http.StatusBadRequest, "invalid token input", err.Error())
			return
		case errors.Is(err, service.ErrorIdNotFound):
			helper.ToResponseJson(ctx, http.StatusUnauthorized, "unauthorized", err.Error())
			return
		default:
			helper.ToResponseJson(ctx, http.StatusInternalServerError, "internal server error", err.Error())
			return
		}
	}

	helper.ToResponseJson(ctx, http.StatusOK, "success", result)
}
