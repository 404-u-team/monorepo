package net

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Login    string `json:"login"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type AuthUserRequest struct {
	Email    *string `json:"email"`
	Login    *string `json:"login"`
	Password string  `json:"password"`
}
