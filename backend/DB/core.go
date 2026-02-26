package db

import (
	system "DevSpace/System"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBManager struct {
	Config     *system.Config
	Connection *gorm.DB
}

func connectDB(loggerEnabled bool, dsn string) (*gorm.DB, error) {
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

		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger,
		})

	} else {
		return gorm.Open(postgres.Open(dsn))
	}
}

func InitDB(c *system.Config) (DBManager, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)

	db, dbErr := connectDB(c.EnableGORMLogger, dsn)

	if dbErr != nil {
		return DBManager{}, dbErr
	}

	return DBManager{Config: c, Connection: db}, nil
}

func (m *DBManager) AutoMigrate() error {
	return m.Connection.AutoMigrate(&User{}, &Idea{}, &Notification{}, &Project{}, &Skill_Category{}, &User_Skill{}, &Chat{}, &Project_Slot{}, &Request{}, &Message{}, &Chat_Member{})
}

func (m *DBManager) InsertNewEntity(e DBEntety) {
	e.WriteToDB(m.Connection)
}

func (m *DBManager) SelectFieldsAll(e DBEntety, fields ...string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	res := m.Connection.Table(e.TableName()).Select(fields).Find(&results)
	if res.Error != nil {
		return nil, res.Error
	}

	return results, nil
}
