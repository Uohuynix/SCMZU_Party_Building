package config

import (
	"fmt"
	"os"
)

// Config holds application configuration read from env or defaults.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	AppPort string
}

// LoadConfigFromEnv loads configuration from environment variables,
// with sensible defaults for easy local testing.
func LoadConfigFromEnv() *Config {
	c := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "Hjh994600!"),
		DBName:     getEnv("DB_NAME", "party_school"),
		AppPort:    getEnv("SERVER_PORT", "8080"),
	}
	return c
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// MySQLDSNWithMultiStatements returns DSN enabling multiStatements
func (c *Config) MySQLDSNWithMultiStatements() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

// MySQLDSN returns DSN for GORM
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
