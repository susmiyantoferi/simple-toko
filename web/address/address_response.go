package web

import "time"

type UserInfo struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AddressResponse struct {
	User      UserInfo  `json:"user"`
	ID        uint      `json:"address_id"`
	Addresses string    `json:"addresses"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
