package web

type InventoryCreateRequest struct {
	Location string `validate:"required,min=1,max=100" json:"location"`
}
