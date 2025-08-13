package web

type ProductUpdateRequest struct {
	ID          uint     `validate:"required"`
	InventoryID *uint    `validate:"omitempty" json:"inventory_id,omitempty"`
	Name        *string  `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Price       *float64 `validate:"omitempty" json:"price,omitempty"`
	Description *string  `validate:"omitempty,min=1,max=255" json:"description,omitempty"`
}
