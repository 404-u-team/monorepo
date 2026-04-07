package models

import "github.com/google/uuid"

type UserFavoriteIdea struct {
	UserID uuid.UUID  `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	IdeaID *uuid.UUID `gorm:"column:idea_id;type:uuid;primaryKey" json:"idea_id"`

	User User `gorm:"foreignKey:UserID" json:"-"`
	Idea Idea `gorm:"foreignKey:IdeaID" json:"-"`
}

func (uf *UserFavoriteIdea) TableName() string { return "User_Favorite_Idea" }
