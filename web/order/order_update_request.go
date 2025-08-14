package web

type OrderUpdateRequest struct {
	ID        uint `validate:"required"`
	AddressID uint `validate:"required" json:"address_id"`
}
