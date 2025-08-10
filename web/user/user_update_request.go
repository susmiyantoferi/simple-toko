package web

type UserUpdateRequest struct {
	Name     *string `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Password *string `validate:"omitempty,min=1,max=255" json:"password,omitempty"`
}
