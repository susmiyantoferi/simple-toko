package web

type WebResponse struct {
	Code    int         `json:"Code"`
	Message string      `json:"Message"`
	Data    interface{} `json:"Data,omitempty"`
}
