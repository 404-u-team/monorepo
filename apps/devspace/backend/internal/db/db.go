package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/auth"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newPostgresConnection(loggerEnabled bool, dsn string) *gorm.DB {
	gormConfig := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		TranslateError: true,
	}
	if loggerEnabled {
		gormLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		)

		gormConfig.Logger = gormLogger
		db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			log.Fatalf("got error during connection to postgres, %v", err)
		}

		return db
	}

	db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
	if err != nil {
		log.Fatalf("got error during connection to postgres, %v", err)
	}

	return db
}

func InitDB(c *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)

	db := newPostgresConnection(c.EnableGORMLogger, dsn)
	return db
}

func CreateAdmin(dbConn *gorm.DB, config config.Config) {
	hash, err := auth.HashPassword("admin", &config)
	if err != nil {
		log.Fatalln("error when creating admin user: ", err)
	}

	var count int64
	result := dbConn.Model(&models.User{}).Where("nickname = ?", "admin").Count(&count)
	if result.Error != nil {
		log.Fatalln("error when creating admin user: ", err)
	}
	if count != 0 {
		log.Println("admin user already created")
		return
	}

	result = dbConn.Create(&models.User{Email: "admin@mail.com", PasswordHash: hash, Nickname: "admin", IsAdmin: true})
	if result.Error != nil {
		log.Fatalln("error when creating admin user: ", err)
	}
}
