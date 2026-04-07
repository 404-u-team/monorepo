package models

import "github.com/google/uuid"

type UserFavorite struct {
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	IdeaID    *uuid.UUID `gorm:"column:idea_id;type:uuid;primaryKey" json:"idea_id"`
	ProjectID *uuid.UUID `gorm:"column:project_id;type:uuid;primaryKey" json:"project_id"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Idea    Idea    `gorm:"foreignKey:IdeaID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
}

func (uf *UserFavorite) TableName() string { return "User_Favorite" }
