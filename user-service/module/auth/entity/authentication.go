package entity

type AuthenticationUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthenticationUser struct {
	AccessToken string `json:"access_token"`
}

type AuthenticationUserResponse struct {
	Auth *AuthenticationUser `json:"auth"`
	Meta *Meta               `json:"meta"`
}
