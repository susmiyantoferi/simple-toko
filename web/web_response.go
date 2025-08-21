package web

type WebResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	CurrentPage int64       `json:"current_page,omitempty"`
	TotalPage   int         `json:"total_page,omitempty"`
	TotalItems  int64       `json:"total_items,omitempty"`
	Data        interface{} `json:"data,omitempty"`
}
