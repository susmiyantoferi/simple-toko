package web

type AddressUpdateRequest struct {
	ID        uint   `validate:"required"`
	UserID    uint   `validate:"required"`
	Addresses string `validate:"required" json:"addresses"`
}
