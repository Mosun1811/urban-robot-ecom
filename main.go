package main

import (
	"errors"
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
	cartService := service.CartService{
		Repo:        cartRepo,
		ProductRepo: productRepo,
	}

	orderService := service.OrderService{
		OrderRepo:   orderRepo,
		CartRepo:    cartRepo,
		ProductRepo: productRepo,
	}
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
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)

	// ----------------------------------------------------
	// 1) Seed demo product list ONLY if table is empty
	// ----------------------------------------------------
	if productCount == 0 {
		log.Println("No products found → Seeding demo products...")

		demoProducts := []models.Product{
			{Name: "Apple iPhone 16", Description: "Latest Apple smartphone", Category: "electronics", PriceCents: 99900, Stock: 50},
			{Name: "Samsung Galaxy S25", Description: "Flagship Android phone", Category: "electronics", PriceCents: 89900, Stock: 40},
			{Name: "Nike Air Max", Description: "Comfortable running shoes", Category: "fashion", PriceCents: 12000, Stock: 100},
			{Name: "Adidas Ultraboost", Description: "High performance running shoes", Category: "fashion", PriceCents: 14500, Stock: 80},
			{Name: "Sony WH-2000XM6", Description: "Noise cancelling headphones", Category: "electronics", PriceCents: 35000, Stock: 25},
			{Name: "4K Smart TV", Description: "55-inch Ultra HD smart TV", Category: "electronics", PriceCents: 250000, Stock: 15},
			{Name: "Kitchen Blender", Description: "High power blender", Category: "home", PriceCents: 8000, Stock: 60},
			{Name: "Office Chair", Description: "Ergonomic chair", Category: "furniture", PriceCents: 8500, Stock: 30},
			{Name: "Gaming Laptop", Description: "RTX graphics gaming machine", Category: "electronics", PriceCents: 180000, Stock: 10},
			{Name: "Electric Kettle", Description: "Stainless steel kettle", Category: "home", PriceCents: 3000, Stock: 50},
		}

		for _, p := range demoProducts {
			product := p
			db.Create(&product)

			// Create Matching Stock row
			db.Create(&models.Stock{
				ProductID: product.ID,
				Quantity:  int(product.Stock),
			})
		}

		log.Println("Demo products + stock seeded.")
	}

	// ----------------------------------------------------
	// 2) SELF-HEALING: Ensure every product has Stock row
	// ----------------------------------------------------
	var products []models.Product
	db.Find(&products)

	for _, p := range products {
		var stock models.Stock
		err := db.Where("product_id = ?", p.ID).First(&stock).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Missing stock for product %d → Creating now...\n", p.ID)

			db.Create(&models.Stock{
				ProductID: p.ID,
				Quantity:  int(p.Stock), // Use Product.Stock as fallback
			})
		}
	}

	log.Println("Stock self-healing complete — all products now have stock entries.")
}
