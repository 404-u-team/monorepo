package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Nickname string `json:"nickname" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required"` // email or nickname
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string
	RefreshToken string
}
