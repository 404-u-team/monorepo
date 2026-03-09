package dto

type UpdateUserRequest struct {
	Nickname *string `json:"nickname" binding:"min=3,max=50"`
	Bio      *string `json:"bio" binding:"min=3,max=255"`
}
