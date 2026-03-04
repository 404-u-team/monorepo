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

	// &User{}, &Idea{}, &Notification{}, &Project{}, &Skill_Category{}, &User_Skill{}, &Chat{}, &Project_Slot{}, &Request{}, &Message{}, &Chat_Member{}
	if err := connection.AutoMigrate(&models.User{}); err != nil {
		log.Fatalln("Произошла ошибка при миграции: ", err)
	}
}
