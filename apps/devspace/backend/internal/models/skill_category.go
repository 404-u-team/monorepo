package models

type SkillCategory struct {
	ID       uint   `gorm:"primaryKey; column:id"`
	ParentId uint   `gorm:"column:parent_id"`
	Name     string `gorm:"column:name; not null"`

	Parent *SkillCategory `gorm:"foreignKey:ParentId"`
}

func (skillCategory *SkillCategory) TableName() string { return "Skill_Category" }
