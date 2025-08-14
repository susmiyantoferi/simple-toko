package web

type OrderUpdateStatusRequest struct {
	ID             uint    `validate:"required"`
	StatusOrder    *string `validate:"omitempty,oneof=waiting confirmed canceled" json:"status_order,omitempty"`
	StatusDelivery *string `validate:"omitempty,oneof=on_process delivered canceled" json:"status_delivery,omitempty"`
}
