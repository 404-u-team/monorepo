package models

import "github.com/google/uuid"

type UserFavoriteProject struct {
	UserID    uuid.UUID `gorm:"column:user_id;type:uuid;primaryKey" json:"user_id"`
	ProjectID uuid.UUID `gorm:"column:project_id;type:uuid;primaryKey" json:"project_id"`

	User    User    `gorm:"foreignKey:UserID" json:"-"`
	Project Project `gorm:"foreignKey:ProjectID" json:"-"`
}

func (pf *UserFavoriteProject) TableName() string { return "User_Favorite_Project" }
