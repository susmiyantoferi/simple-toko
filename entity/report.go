package entity

type SalesReport struct {
	Month      string  `json:"month"`
	TotalQty   int     `json:"total_qty"`
	TotalSales float64 `json:"total_sales"`
}

type TopProduct struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalQty    int64   `json:"total_qty"`
	TotalSales  float64 `json:"total_sales"`
}
