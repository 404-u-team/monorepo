package models

type Project_Slot struct {
	Id              uint   `gorm:"primaryKey; column:id"`
	ProjectId       uint   `gorm:"column:project_id; not null"`
	SkillCategoryId uint   `gorm:"column:skill_category_id; not null"`
	UserId          uint   `gorm:"column:user_id"`
	Status          string `gorm:"column:status; not null"`

	Project Project       `gorm:"foreignKey:ProjectId"`
	Skill   SkillCategory `gorm:"foreignKey:SkillCategoryId"`
	User    User          `gorm:"foreignKey:UserId"`
}

func (ps *Project_Slot) TableName() string { return "Project_Slot" }
