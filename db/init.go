package db

import (
	"fmt"
	"log"

	"github.com/matiuskm/go-commerce/config"
	"github.com/matiuskm/go-commerce/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetEnv("DB_HOST", "localhost"),
		config.GetEnv("DB_USER", ""),
		config.GetEnv("DB_PASSWORD", ""),
		config.GetEnv("DB_NAME", "ecommerce"),
		config.GetEnv("DB_PORT", "5432"),
	  )
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Order{},
		&models.OrderItem{},
		&models.Cart{},
		&models.CartItem{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = database
	log.Println("Database connected and migrated")
}