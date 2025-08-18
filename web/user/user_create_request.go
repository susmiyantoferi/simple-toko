package web

type UserCreateRequest struct {
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=100" json:"email"`
	Password string `validate:"required,min=1,max=255" json:"password"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,email,min=1,max=100" json:"email"`
	Password string `validate:"required,min=1,max=255" json:"password"`
}

type UserRefreshTokenRequest struct {
	TokenRefresh string `validate:"required" json:"token_refresh"`
}
