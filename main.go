package main

import (
	"log"
	"net/http"
	"os"

	"futuremarket/db"
	"futuremarket/handlers"
	"futuremarket/routes"
)

func main() {

	database := db.InitDB()

	//    Gonna leave these empty for now.
	//    Later, the team will work through services/repositories into these.
	authHandler := &handlers.AuthHandler{

		DB: database,
	}

	productHandler := &handlers.ProductHandler{
		//   ProductRepo, StockRepo
		DB: database,
	}

	cartHandler := &handlers.CartHandler{
		// which will use CartRepos, Product/StockRepo
		DB: database,
	}

	orderHandler := &handlers.OrderHandler{
		//  which will handle checkout + transactions
		DB: database,
	}

	reviewHandler := &handlers.ReviewHandler{
		//   ReviewRepo + ProductRepo
		DB: database,
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
