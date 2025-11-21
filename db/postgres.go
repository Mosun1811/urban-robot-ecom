package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"futuremarket/models"
)


var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Println("warning: could not load .env file, falling back to environment variables")
	}

	connStr := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// 3. Open the database using GORM and the Postgres driver.
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("connected to database successfully!")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Stock{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.OrderItem{},
		&models.Review{},
	)
	if err != nil {
		log.Fatalf("unable to migrate schema: %v", err)
	}

	return DB
}
