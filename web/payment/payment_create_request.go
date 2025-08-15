package web

type PaymentCreateRequest struct {
	OrderID uint   `form:"order_id" binding:"required"`
	Image   string `form:"-"`
}
