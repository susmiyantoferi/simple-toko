package web

type ProductStockUpdateRequest struct {
	ID    uint `validate:"required" json:"id"`
	Stock int  `validate:"required,gt=0" json:"stock"`
}
