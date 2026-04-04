package dto

// body запроса на создание отклика на проект
type CreateProjectRequestApplyRequest struct {
	CoverLetter *string `json:"cover_letter" binding:"omitempty,min=1,max=255"`
}

// body запроса на создание запроса на приглашение в проект
type CreateProjectRequestInviteRequest struct {
	UserID      string  `json:"user_id" binding:"required"`
	CoverLetter *string `json:"cover_letter" binding:"omitempty,min=1,max=255"`
}
