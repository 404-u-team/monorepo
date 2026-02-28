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
	"gorm.io/gorm/schema"
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
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
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

func (m *DBManager) InsertNewEntity(e DBEntety) error {
	return e.WriteToDB(m.Connection)
}

// берет поля и структуру. Возвращает хэштаблицы со связкой поле-значение
func (m *DBManager) SelectFieldsAll(e DBEntety, fields ...string) ([]map[string]any, error) {
	var results []map[string]any
	res := m.Connection.Table(e.TableName()).Select(fields).Find(&results)
	if res.Error != nil {
		return nil, res.Error
	}

	return results, nil
}

func (m *DBManager) DeleteByID(e DBEntety, id uint) (bool, error) {
	result := m.Connection.Delete(e, id)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}
