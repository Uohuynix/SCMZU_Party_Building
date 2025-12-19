package configs

import (
	"os"
)

// Config application configuration structure
type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	ServerPort  string
	LogLevel    string
	Environment string
}

// GetConfig loads configuration from environment variables
func GetConfig() *Config {
	return &Config{
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "Hjh994600!"),
		DBName:      getEnv("DB_NAME", "party_school"),
		JWTSecret:   getEnv("JWT_SECRET", "11223344"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
