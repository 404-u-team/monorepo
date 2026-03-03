package db

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

func Migrate(connection *gorm.DB) {
	// &User{}, &Idea{}, &Notification{}, &Project{}, &Skill_Category{}, &User_Skill{}, &Chat{}, &Project_Slot{}, &Request{}, &Message{}, &Chat_Member{}
	err := connection.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Произошла ошибка при миграции, %v", err)
	}
}
