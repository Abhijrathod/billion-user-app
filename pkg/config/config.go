package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
// Services will only use the parts they need.
type Config struct {
	// --- Generic Service Config ---
	AppEnv string // "development", "staging", "production"
	Port   string

	// --- Database (Postgres) ---
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSslMode  string

	// --- Caching (Redis) ---
	RedisAddress string

	// --- Messaging (Kafka) ---
	KafkaBrokers string

	// --- Auth (JWT) ---
	JWTSecret string
}

// LoadConfig loads configuration from environment variables
// It will load from a .env file if one is present in the service's directory.
func LoadConfig(envPath ...string) (*Config, error) {
	// Look for .env file.
	// We allow passing a path (e.g., "../.env") for different service depths
	// If no path is provided, it checks the current directory.
	if len(envPath) > 0 {
		err := godotenv.Load(envPath[0])
		if err != nil {
			log.Printf("Warning: could not load .env file from %s. Using environment variables.", envPath[0])
		}
	} else {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: no .env file found. Using environment variables.")
		}
	}

	// Read values from environment
	cfg := &Config{
		AppEnv:       getEnv("APP_ENV", "development"),
		Port:         getEnv("PORT", "3000"),
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "admin"),
		DBPassword:   getEnv("DB_PASSWORD", "secret"),
		DBName:       getEnv("DB_NAME", "postgres"),
		DBSslMode:    getEnv("DB_SSLMODE", "disable"),
		RedisAddress: getEnv("REDIS_ADDRESS", "localhost:6379"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		JWTSecret:    getEnv("JWT_SECRET", "super-secret-key"),
	}

	return cfg, nil
}

// getEnv is a helper to read an environment variable or return a default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
