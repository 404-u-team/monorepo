package models

type UserSkill struct {
	UserID  uint `gorm:"column:user_id; primaryKey"`
	SkillId uint `gorm:"column:skill_id; primaryKey"`

	User  User          `gorm:"foreignKey:UserId"`
	Skill SkillCategory `gorm:"foreignKey:SkillId"`
}

func (us *UserSkill) TableName() string { return "User_Skill" }
