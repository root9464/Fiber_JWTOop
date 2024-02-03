package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var slowLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags),
	logger.Config{
		SlowThreshold: 1 * time.Microsecond,
		LogLevel:      logger.Silent,
		Colorful:      true,
	},
)

type Database interface {
	Connect() (*gorm.DB, error) 
	Migration(data interface{}) error
	CreateTestData(c *fiber.Ctx) error
}

type Data struct {
	DB *gorm.DB
}

func (method *Data) Connect() (*gorm.DB, error) {
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("Error loading .env file: %w", err)
    }
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
        os.Getenv("PGHOST"),
        os.Getenv("PGUSER"),
        os.Getenv("PGPASSWORD"),
        os.Getenv("PGDATABASE"),
        os.Getenv("PGPORT"),
    )
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: slowLogger})
    if err != nil {
        return nil, fmt.Errorf("Failed to connect to the database: %w", err)
    }
    method.DB = db
    log.Println("Database connection successful.")
    return db, nil
}
