package web

type PaymentUpdateRequest struct {
	OrderID uint    `validate:"required" json:"order_id"`
	Status  *string `validate:"omitempty,oneof=waiting confirmed canceled" json:"status,omitempty"`
	Image   *string `validate:"omitempty" json:"image,omitempty"`
}
