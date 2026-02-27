package system

import (
	consts "DevSpace/Consts"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// API
	APIPort uint

	// Database
	DBPort     uint
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	// Features
	EnableGORMLogger bool
	DebugMode        bool

	// Argon
	Memory      uint
	Iterations  uint
	Parallelism uint
	SaltLength  uint
	KeyLength   uint
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден... Использую системное окружение")
	}

	cfg := &Config{
		// API с дефолтом из consts
		APIPort: getEnvAsUint("API_PORT", consts.APIPortDefault),

		// DB с дефолтами
		DBPort:     getEnvAsUint("DB_PORT", 5432),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "devspace"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),

		// Фичи
		EnableGORMLogger: getEnvAsBool("GORM_LOGGER", false),
		DebugMode:        getEnvAsBool("DEBUG", false),

		// Аргон
		Memory:      getEnvAsUint("ARGON_MEMORY_USAGE", 54),
		Iterations:  getEnvAsUint("ARGON_ITERATIONS", 5),
		Parallelism: getEnvAsUint("ARGON_PARALL", 5),
		SaltLength:  getEnvAsUint("ARGON_SALT_LEN", 16),
		KeyLength:   getEnvAsUint("ARGON_KEY_LEN", 32),
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsUint(key string, fallback uint) uint {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.ParseUint(value, 10, 32); err == nil {
			return uint(intVal)
		}
	}
	return fallback
}

func getEnvAsBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return fallback
}
