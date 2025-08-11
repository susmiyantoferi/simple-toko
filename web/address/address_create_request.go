package web

type AddressCreateRequest struct {
	UserID uint `validate:"required" json:"user_id"`
	Addresses string `validate:"required" json:"addresses"`
}
