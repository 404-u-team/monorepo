package db

import (
	"time"

	"gorm.io/gorm"
)

type DBEntety interface {
	WriteToDB(*gorm.DB)
	TableName() string
}

type User struct {
	ID           uint   `gorm:"primaryKey; column:id"`
	Email        string `gorm:"unique; column:email; not null"`
	PasswordHash string `gorm:"column:password_hash; not null"`
	Nickname     string `gorm:"column:nickname; not null"`
	AvatarUrl    string `gorm:"column:avatar_url; not null"`
	Status       string `gorm:"column:status; not null"`
	Bio          string `gorm:"column:bio; not null"`
}

func (u *User) WriteToDB(d *gorm.DB) {
	d.Create(&u)
}

func (u *User) TableName() string { return "User" }

type Notification struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	UserId    uint      `gorm:"column:user_id; not null"`
	Message   string    `gorm:"column:message; not null"`
	IsRead    bool      `gorm:"column:is_read; not null"`
	CreatedAt time.Time `gorm:"column:created_at; not null"`

	User User `gorm:"foreignKey:UserId"`
}

func (n *Notification) WriteToDB(d *gorm.DB) {
	d.Create(&n)
}

func (n *Notification) TableName() string { return "Notification" }

type Idea struct {
	Id             uint   `gorm:"primaryKey; column:id"`
	AuthorId       uint   `gorm:"column:author_id"`
	Title          string `gorm:"column:title; not null"`
	Description    string `gorm:"column:description; not null"`
	ViewsCount     uint   `gorm:"column:views_count; not null"`
	FavoritesCount uint   `gorm:"column:favorites_count; not null"`

	Author User `gorm:"foreignKey:AuthorId"`
}

func (i *Idea) WriteToDB(d *gorm.DB) {
	d.Create(&i)
}

func (i *Idea) TableName() string { return "Idea" }

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

func (p *Project) WriteToDB(d *gorm.DB) {
	d.Create(&p)
}
func (p *Project) TableName() string { return "Project" }

type Skill_Category struct {
	Id       uint   `gorm:"primaryKey; column:id"`
	ParentId uint   `gorm:"column:parent_id"`
	Name     string `gorm:"column:name; not null"`

	Parent *Skill_Category `gorm:"foreignKey:ParentId"`
}

func (sc *Skill_Category) WriteToDB(d *gorm.DB) {
	d.Create(&sc)
}

func (sc *Skill_Category) TableName() string { return "Skill_Category" }

type User_Skill struct {
	UserId  uint `gorm:"column:user_id; primaryKey"`
	SkillId uint `gorm:"column:skill_id; primaryKey"`

	User  User           `gorm:"foreignKey:UserId"`
	Skill Skill_Category `gorm:"foreignKey:SkillId"`
}

func (us *User_Skill) WriteToDB(d *gorm.DB) {
	d.Create(&us)
}

func (us *User_Skill) TableName() string { return "User_Skill" }

type Chat struct {
	Id        uint      `gorm:"primaryKey; column:id"`
	ProjectId uint      `gorm:"column:project_id"`
	Title     string    `gorm:"column:title"`
	Type      string    `gorm:"column:type; not null"`
	CreatedAt time.Time `gorm:"column:created_at"`

	Project Project `gorm:"foreignKey:ProjectId"`
}

func (c *Chat) WriteToDB(d *gorm.DB) {
	d.Create(&c)
}
func (c *Chat) TableName() string { return "Chat" }

type Chat_Member struct {
	ChatId     uint      `gorm:"column:chat_id; primaryKey"`
	UserId     uint      `gorm:"column:user_id; primaryKey"`
	JoinedAt   time.Time `gorm:"column:joined_at; not null"`
	LastReadAt time.Time `gorm:"column:last_read_at; not null"`

	Chat Chat `gorm:"foreignKey:ChatId"`
	User User `gorm:"foreignKey:UserId"`
}

func (cm *Chat_Member) WriteToDB(d *gorm.DB) {
	d.Create(&cm)
}

func (cm *Chat_Member) TableName() string { return "Chat_Member" }

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

func (ps *Project_Slot) WriteToDB(d *gorm.DB) {
	d.Create(&ps)
}

func (ps *Project_Slot) TableName() string { return "Project_Slot" }

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

func (r *Request) WriteToDB(d *gorm.DB) {
	d.Create(&r)
}
func (r *Request) TableName() string { return "Request" }

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

func (m *Message) WriteToDB(d *gorm.DB) {
	d.Create(&m)
}

func (m *Message) TableName() string { return "Message" }
