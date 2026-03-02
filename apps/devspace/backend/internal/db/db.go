package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"github.com/404-u-team/monorepo/apps/devspace/backend/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func newPostgresConnection(loggerEnabled bool, dsn string) *gorm.DB {
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

		gormConfig := gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: gormLogger,
		}

		db, err := gorm.Open(postgres.Open(dsn), &gormConfig)
		if err != nil {
			log.Fatalf("got error during connection to postgres, %v", err)
		}

		return db
	}

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("got error during connection to postgres, %v", err)
	}

	return db
}

func InitDB(c *config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)

	db := newPostgresConnection(c.EnableGORMLogger, dsn)
	return db
}

func CreateEntity(db *gorm.DB, entity any) error {
	res := db.Create(entity)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return UniqueKeyDuplErr
		}
		log.Println("failed to create %v: %v", entity, res.Error)
	}

	if res.RowsAffected == 0 {
		return UniqueKeyDuplErr
	}

	return nil
}
