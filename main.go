package main

import (
	"log"
	"net/http"
	"os"

	"futuremarket/db"
	"futuremarket/handlers"
	"futuremarket/models"
	"futuremarket/repository"
	"futuremarket/routes"
	"futuremarket/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {

	database := db.InitDB()
	seedAdminUser(database)
	seedDemoProducts(database)

	// ----------------------------
	// REPOSITORIES
	// ----------------------------
	userRepo := repository.UserRepo{DB: database}
	cartRepo := repository.CartRepo{DB: database}
	orderRepo := repository.OrderRepo{DB: database}
	productRepo := repository.ProductRepo{DB: database}
	reviewRepo := repository.ReviewRepo{DB: database}
	blacklistRepo := repository.NewBlacklistRepository(database) // Blacklist

	// ----------------------------
	// SERVICES
	// ----------------------------
	userService := service.UserService{Repo: userRepo}
	cartService := service.CartService{Repo: cartRepo}
	orderService := service.OrderService{Repo: orderRepo}
	productService := service.ProductService{Repo: productRepo}
	reviewService := service.ReviewService{Repo: reviewRepo}
	blacklistService := service.BlacklistService{Repo: blacklistRepo} // ⭐ NEW

	// ----------------------------
	// HANDLERS
	// ----------------------------

	authHandler := &handlers.AuthHandler{
		Service:          userService,
		BlacklistService: blacklistService,
	}

	productHandler := &handlers.ProductHandler{
		Service: productService,
	}

	cartHandler := &handlers.CartHandler{
		Service: cartService,
	}

	orderHandler := &handlers.OrderHandler{
		Service: orderService,
	}

	reviewHandler := &handlers.ReviewHandler{
		Service: reviewService,
	}

	// ----------------------------
	// ROUTER
	// ----------------------------
	router := routes.SetupRouter(
		authHandler,
		productHandler,
		cartHandler,
		orderHandler,
		reviewHandler,
		blacklistService,
	)

	// ----------------------------
	// START SERVER
	// ----------------------------
	addr := os.Getenv("APP_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("FutureMarket API starting on... %s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// ============================================================
// SEED ADMIN USER — Must be defined OUTSIDE main()
// ============================================================
//

func seedAdminUser(db *gorm.DB) {
	var admin models.User

	// Check if admin already exists
	err := db.Where("email = ?", "admin@futuremarket.com").First(&admin).Error
	if err == nil {
		return // admin exists already
	}

	// Create admin password
	password := "AdminPass123!"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Create admin user
	admin = models.User{
		Name:         "System Admin",
		Email:        "admin@futuremarket.com",
		PasswordHash: string(hashed),
		Role:         "admin",
	}

	db.Create(&admin)

	log.Println("Admin user created: admin@futuremarket.com / AdminPass123!")
}

func seedDemoProducts(db *gorm.DB) {
	var count int64
	db.Model(&models.Product{}).Count(&count)

	if count > 0 {
		log.Println("Demo products already exist, skipping seed.")
		return
	}

	demoProducts := []models.Product{
		{Name: "Apple iPhone 16", Description: "Latest Apple smartphone", Category: "electronics", PriceCents: 99900, Stock: 50},
		{Name: "Samsung Galaxy S25", Description: "Flagship Android phone", Category: "electronics", PriceCents: 89900, Stock: 40},
		{Name: "Nike Air Max", Description: "Comfortable running shoes", Category: "fashion", PriceCents: 12000, Stock: 100},
		{Name: "Adidas Ultraboost", Description: "High performance running shoes", Category: "fashion", PriceCents: 14500, Stock: 80},
		{Name: "Sony Headphones WH-2000XM6", Description: "Noise cancelling headphones", Category: "electronics", PriceCents: 35000, Stock: 25},
		{Name: "4K Smart TV", Description: "55-inch Ultra HD smart TV", Category: "electronics", PriceCents: 250000, Stock: 15},
		{Name: "Kitchen Blender", Description: "High power blender", Category: "home", PriceCents: 8000, Stock: 60},
		{Name: "Office Chair", Description: "Ergonomic chair", Category: "furniture", PriceCents: 8500, Stock: 30},
		{Name: "Gaming Laptop", Description: "RTX graphics gaming machine", Category: "electronics", PriceCents: 180000, Stock: 10},
		{Name: "Electric Kettle", Description: "Stainless steel kettle", Category: "home", PriceCents: 3000, Stock: 50},
	}

	for _, product := range demoProducts {
		db.Create(&product)
	}

	log.Println("Seeded 10 demo products successfully!")
}
