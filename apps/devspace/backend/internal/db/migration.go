package db

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

func Migrate(connection *gorm.DB) {
	err := connection.AutoMigrate(&models.User{}, &models.Idea{}, &models.Notification{}, &models.Project{}, &models.SkillCategory{}, &models.UserSkill{}, &models.Chat{}, &models.ProjectSlot{}, &models.Request{}, &models.Message{}, &models.ChatMember{})
	if err != nil {
		log.Fatalf("Произошла ошибка при миграции, %v", err)
	}
}
