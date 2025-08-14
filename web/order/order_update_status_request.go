package web

type OrderUpdateStatusRequest struct {
	ID             uint   `validate:"required"`
	StatusOrder    string `validate:"required,oneof=waiting confirmed canceled" json:"status_order"`
	StatusDelivery string `validate:"required,oneof=on_process delivered canceled" json:"status_delivery"`
}
