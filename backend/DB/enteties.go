package db

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey; column:id"`
	Email        string `gorm:"unique; column:email; not null"`
	PasswordHash string `gorm:"column:password_hash; not null"`
	Nickname     string `gorm:"column:nickname; not null"`
	AvatarUrl    string `gorm:"column:avatar_url; not null"`
	Status       string `gorm:"column:status; not null"`
	Bio          string `gorm:"column:bio; not null"`
}

type Notification struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	UserId    uint      `gorm:"column:user_id; not null"`
	Message   string    `gorm:"column:message; not null"`
	IsRead    bool      `gorm:"column:is_read; not null"`
	CreatedAt time.Time `gorm:"column:created_at; not null"`

	User User `gorm:"foreignKey:UserId"`
}

type Idea struct {
	Id             uint   `gorm:"primaryKey; column:id"`
	AuthorId       uint   `gorm:"column:author_id"`
	Title          string `gorm:"column:title; not null"`
	Description    string `gorm:"column:description; not null"`
	ViewsCount     uint   `gorm:"column:views_count; not null"`
	FavoritesCount uint   `gorm:"column:favorites_count; not null"`

	Author User `gorm:"foreignKey:AuthorId"`
}

type Project struct {
	Id          uint      `gorm:"primaryKey; column:id"`
	LeaderId    uint      `gorm:"column:leader_id; not null"`
	IdeaId      uint      `gorm:"column:idea_id"`
	Title       string    `gorm:"column:title; not null"`
	Description string    `gorm:"column:descriprion; not null"`
	Status      string    `gorm:"column:status; not null"`
	CreatedAt   time.Time `gorm:"column:created_at; not null"`

	Leader User `gorm:"foreignKey:LeaderId"`
	Idea   Idea `gorm:"foreignKey:IdeaId"`
}

type Skill_Category struct {
	Id       uint   `gorm:"primaryKey; column:id"`
	ParentId uint   `gorm:"column:parent_id"`
	Name     string `gorm:"column:name; not null"`

	Parent *Skill_Category `gorm:"foreignKey:ParentId"`
}

type User_Skill struct {
	UserId  uint `gorm:"column:user_id; primaryKey"`
	SkillId uint `gorm:"column:skill_id; primaryKey"`

	User  User           `gorm:"foreignKey:UserId"`
	Skill Skill_Category `gorm:"foreignKey:SkillId"`
}

type Chat struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	ProjectId uint      `gorm:"column:project_id"`
	Title     string    `gorm:"column:title"`
	Type      string    `gorm:"column:type; not null"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Project Project `gorm:"foreignKey:ProjectId"`
}

type Chat_Member struct {
	ChatId     uint      `gorm:"column:chat_id; primaryKey"`
	UserId     uint      `gorm:"column:user_id; primaryKey"`
	JoinedAt   time.Time `gorm:"column:joined_at; not null"`
	LastReadAt time.Time `gorm:"column:last_read_at; not null"`

	Chat Chat `gorm:"foreignKey:ChatId"`
	User User `gorm:"foreignKey:UserId"`
}

type Project_Slot struct {
	Id              uint   `gorm:"primaryKey; column:id"`
	ProjectId       uint   `gorm:"column:project_id; not null"`
	SkillCategoryId uint   `gorm:"column:skill_category_id; not null"`
	UserId          uint   `gorm:"column:user_id"`
	Status          string `gorm:"column:status; not null"`

	Project Project        `gorm:"foreignKey:ProjectId"`
	Skill   Skill_Category `gorm:"foreignKey:SkillCategoryId"`
	User    User           `gorm:"foreignKey:UserId"`
}

type Request struct {
	Id          uint   `gorm:"primaryKey; column:id"`
	SlotId      uint   `gorm:"column:slot_id; not null"`
	UserId      uint   `gorm:"column:user_id; not null"`
	Type        string `gorm:"column:type; not null"`
	Status      string `gorm:"column:status; not null"`
	CoverLetter string `gorm:"column:cover_letter; not null"`

	Slot Project_Slot `gorm:"foreignKey:SlotId"`
	User User         `gorm:"foreignKey:UserId"`
}

type Message struct {
	Id       uint      `gorm:"primaryKey; column:id"`
	ChatId   uint      `gorm:"column:chat_id; not null"`
	SenderId uint      `gorm:"column:sender_id; not null"`
	Content  string    `gorm:"column:content; not null"`
	SentAt   time.Time `gorm:"column:sent_at; not null"`
	IsEdited bool      `gorm:"column:is_edited; not null"`

	Chat   Chat `gorm:"foreignKey:ChatId"`
	Sender User `gorm:"foreignKey:SenderId"`
}
