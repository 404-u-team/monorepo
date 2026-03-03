package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"` // email or nickname
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}
