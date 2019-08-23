package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config is the configuration of the server.
type Config struct {
	AppEnv   string
	TLS      bool
	DBName   string
	Port     int
	KeyFile  string
	CertFile string
	LogLevel int
}

// New creates an instance of config.
func New() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		TLS:      getEnvAsBool("TLS", false),
		DBName:   getEnv("DB_NAME", "prototodos"),
		KeyFile:  getEnv("KEY_FILE", ""),
		CertFile: getEnv("CERT_FILE", ""),
		Port:     getEnvAsInt("PORT", 10000),
		LogLevel: getEnvAsInt("LOG_LEVEL", 0),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper to read an environment variable into a bool or return a default value
func getEnvAsBool(name string, defaultVal bool) bool {
	valStr := getEnv(name, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

// Simple helper function to read an environment variable into an integer
// or return a default value
func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
