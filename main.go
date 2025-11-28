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
	userRepo:= repository.UserRepo{DB: database}
	cartRepo:= repository.CartRepo{DB: database}
	orderRepo:= repository.OrderRepo{DB: database}
	productRepo:= repository.ProductRepo{DB: database}
	reviewRepo:= repository.ReviewRepo{DB: database}

	userService:= service.UserService{Repo: userRepo}
	cartService:= service.CartService{Repo: cartRepo}
	orderService:= service.OrderService{Repo: orderRepo}
	productService:= service.ProductService{Repo: productRepo}
	reviewService:= service.ReviewService{Repo: reviewRepo}
	//    Gonna leave these empty for now.
	//    Later, the team will work through services/repositories into these.
	authHandler := &handlers.AuthHandler{

		Service: userService,
	}

	productHandler := &handlers.ProductHandler{
		//   ProductRepo, StockRepo
		Service: productService,
	}

	cartHandler := &handlers.CartHandler{
		// which will use CartRepos, Product/StockRepo
		Service: cartService,
	}

	orderHandler := &handlers.OrderHandler{
		//  which will handle checkout + transactions
		Service: orderService,
	}

	reviewHandler := &handlers.ReviewHandler{
		//   ReviewRepo + ProductRepo
		Service: reviewService,
	}

	//  Register all routes and middleware.
	router := routes.SetupRouter(
		authHandler,
		productHandler,
		cartHandler,
		orderHandler,
		reviewHandler,
	)

	// Start the HTTP server.
	addr := os.Getenv("APP_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("FutureMarket API starting on...%s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
