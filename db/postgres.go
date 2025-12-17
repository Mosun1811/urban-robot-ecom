package db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"futuremarket/models"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// Load .env ONLY for local development
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set — cannot connect to database")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB = db

	log.Println("Connected to database successfully!")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Stock{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
		&models.TokenBlacklist{},
	)
	if err != nil {
		log.Fatalf("unable to migrate schema: %v", err)
	}

	return DB
}
