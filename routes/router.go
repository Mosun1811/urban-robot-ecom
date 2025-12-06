package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"futuremarket/handlers"
	"futuremarket/middleware"
	"futuremarket/service"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
	reviewHandler *handlers.ReviewHandler,
	blacklistService service.BlacklistService,
) *mux.Router {

	r := mux.NewRouter()

	// Health Check
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the FutureMarket API"))
	})

	// PUBLIC AUTH ROUTES
	r.HandleFunc("/api/v1/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", authHandler.Login).Methods(http.MethodPost)

	// PUBLIC PRODUCT ROUTES
	r.HandleFunc("/api/v1/products", productHandler.ListProducts).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/products/{id}", productHandler.GetProductByID).Methods(http.MethodGet)

	// PUBLIC REVIEWS
	r.HandleFunc("/api/v1/products/{id}/reviews", reviewHandler.ListReviews).Methods(http.MethodGet)

	// ---------------------------------------
	// PROTECTED ROUTES (NEED AUTH TOKEN)
	// ---------------------------------------
	protected := r.PathPrefix("/api/v1").Subrouter()

	protected.Use(
		middleware.AuthMiddlewareConfig{
			BlacklistService: blacklistService,
		}.AuthMiddleware,
	)

	// LOGOUT
	protected.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost)

	// CART
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods(http.MethodGet)
	protected.HandleFunc("/cart", cartHandler.AddToCart).Methods(http.MethodPost)
	protected.HandleFunc("/cart/{product_id}", cartHandler.UpdateCartItem).Methods(http.MethodPatch)
	protected.HandleFunc("/cart/{product_id}", cartHandler.RemoveCartItem).Methods(http.MethodDelete)

	// ORDERS
	protected.HandleFunc("/checkout", orderHandler.Checkout).Methods(http.MethodPost)
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods(http.MethodGet)

	// AUTHENTICATED REVIEW ROUTES
	protected.HandleFunc("/products/{id}/reviews", reviewHandler.CreateOrUpdateReview).Methods(http.MethodPost)

	// ---------------------------------------
	// ADMIN ROUTES
	// ---------------------------------------
	admin := protected.PathPrefix("/admin").Subrouter()

	// IMPORTANT: AdminMiddleware added AFTER AuthMiddleware
	admin.Use(middleware.AdminMiddleware)

	admin.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	admin.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods(http.MethodPatch)

	return r
}
