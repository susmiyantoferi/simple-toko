package web

type ProductItem struct {
	ProductID uint `validate:"required" json:"product_id"`
	Qty      int  `validate:"required,gt=0" json:"qty"`
}

type OrderCreateRequest struct {
	UserID        uint          `validate:"required" json:"user_id"`
	AddressID     uint          `validate:"required" json:"address_id"`
	OrderProducts []ProductItem `validate:"required" json:"order_products"`
}
