package main

import (
	"log"
	"net/http"
	"os"

	"futuremarket/db"
	"futuremarket/handlers"
	"futuremarket/repository"
	"futuremarket/routes"
	"futuremarket/service"
)

func main() {

	database := db.InitDB()

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
