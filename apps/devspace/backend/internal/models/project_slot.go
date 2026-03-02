package models

type ProjectSlot struct {
	ID              uint   `gorm:"primaryKey; column:id"`
	ProjectId       uint   `gorm:"column:project_id; not null"`
	SkillCategoryId uint   `gorm:"column:skill_category_id; not null"`
	UserId          uint   `gorm:"column:user_id"`
	Status          string `gorm:"column:status; not null"`

	Project Project       `gorm:"foreignKey:ProjectId"`
	Skill   SkillCategory `gorm:"foreignKey:SkillCategoryId"`
	User    User          `gorm:"foreignKey:UserId"`
}

func (ps *ProjectSlot) TableName() string { return "Project_Slot" }
