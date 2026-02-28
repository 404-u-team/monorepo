package net

type RegisterUserRequest struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type AuthUserRequest struct {
	Email    *string `json:"email"`
	Nickname *string `json:"nickname"`
	Password string  `json:"password"`
}
