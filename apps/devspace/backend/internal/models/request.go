package models

type Request struct {
	ID          uint   `gorm:"primaryKey; column:id"`
	SlotId      uint   `gorm:"column:slot_id; not null"`
	UserId      uint   `gorm:"column:user_id; not null"`
	Type        string `gorm:"column:type; not null"`
	Status      string `gorm:"column:status; not null"`
	CoverLetter string `gorm:"column:cover_letter; not null"`

	Slot ProjectSlot `gorm:"foreignKey:SlotId"`
	User User        `gorm:"foreignKey:UserId"`
}

func (r *Request) TableName() string { return "Request" }
