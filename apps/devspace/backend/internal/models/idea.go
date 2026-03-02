package models

type Idea struct {
	Id             uint   `gorm:"primaryKey; column:id"`
	AuthorId       uint   `gorm:"column:author_id"`
	Title          string `gorm:"column:title; not null"`
	Description    string `gorm:"column:description; not null"`
	ViewsCount     uint   `gorm:"column:views_count; not null"`
	FavoritesCount uint   `gorm:"column:favorites_count; not null"`

	Author User `gorm:"foreignKey:AuthorId"`
}

func (i *Idea) TableName() string { return "Idea" }
