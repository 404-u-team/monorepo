package db

import (
	"log"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	"gorm.io/gorm"
)

func enablePgCrypto(connection *gorm.DB) error {
	// Включение расширения pgcrypto
	result := connection.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")
	return result.Error
}

func Migrate(connection *gorm.DB) {
	if err := enablePgCrypto(connection); err != nil {
		log.Fatalln("Произошла ошибка при добавления расширения в момент миграции: ", err)
	}

	if err := connection.AutoMigrate(&models.User{}, &models.Idea{}, &models.Notification{}, &models.Project{}, &models.SkillCategory{}, &models.UserSkill{}, &models.Chat{}, &models.ProjectSlot{}, &models.ProjectRequest{}, &models.Message{}, &models.ChatMember{}, &models.UserFavoriteIdea{}); err != nil {
		log.Fatalln("Произошла ошибка при миграции: ", err)
	}
}
