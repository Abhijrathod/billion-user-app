package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// Import your new config package
	"github.com/my-username/billion-user-app/pkg/config"
)

// ConnectDB creates a new GORM database connection
// Each service will call this to get its own connection pool.
func ConnectDB(cfg *config.Config, serviceDBName string) (*gorm.DB, error) {
	// Use the specific DB name for the service, but credentials from the root config
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		serviceDBName, // e.g., "auth_db", "user_db"
		cfg.DBPort,
		cfg.DBSslMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database %s: %w", serviceDBName, err)
	}

	log.Printf("Successfully connected to database: %s", serviceDBName)
	return db, nil
}
