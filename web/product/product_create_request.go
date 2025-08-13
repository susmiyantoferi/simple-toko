package web

type ProductCreateRequest struct {
	InventoryID uint    `validate:"required" json:"inventory_id"`
	Name        string  `validate:"required,min=1,max=100" json:"name"`
	Price       float64 `validate:"required" json:"price"`
	Stock       int     `validate:"required,gt=0" json:"stock"`
	Description string  `validate:"required,min=1,max=225" json:"description"`
}
