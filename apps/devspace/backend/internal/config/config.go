package config

import (
	"log"
	"os"
	"strconv"

	"github.com/lpernett/godotenv"
)

type Config struct {
	APIPort string

	DBPort     string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string

	EnableGORMLogger bool
	DebugMode        bool

	Memory      int
	Iterations  int
	Parallelism int
	SaltLength  int
	KeyLength   int

	JWTSecret              string
	JWTExpirationInSeconds int
}

func LoadConfig() Config {
	err := godotenv.Load("/app/services/auth-service/.env")
	if err != nil {
		log.Printf("got error when tried to load env variables, %v", err)
	}

	return Config{
		APIPort: getEnv("API_PORT", "8080"),

		DBPort:     getEnv("DB_PORT", "5432"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "devspace"),
		SSLMode:    getEnv("DB_SSLMODE", "disable"),

		EnableGORMLogger: getEnvAsBool("GORM_LOGGER", false),
		DebugMode:        getEnvAsBool("DEBUG", false),

		Memory:      getEnvAsInt("ARGON_MEMORY_USAGE", 54),
		Iterations:  getEnvAsInt("ARGON_ITERATIONS", 5),
		Parallelism: getEnvAsInt("ARGON_PARALL", 5),
		SaltLength:  getEnvAsInt("ARGON_SALT_LEN", 16),
		KeyLength:   getEnvAsInt("ARGON_KEY_LEN", 32),

		JWTSecret:              getEnv("JWT_SECRET", "not-secret-anymore"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRE_TIME", 900),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if value, err := strconv.Atoi(value); err == nil {
			return value
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
