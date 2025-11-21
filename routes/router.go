package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"futuremarket/handlers"
	"futuremarket/middleware"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	productHandler *handlers.ProductHandler,
	cartHandler *handlers.CartHandler,
	orderHandler *handlers.OrderHandler,
	reviewHandler *handlers.ReviewHandler,
) *mux.Router {
	r := mux.NewRouter()

	//  Health check / root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Welcome to the FutureMarket API"))
	}).Methods(http.MethodGet)

	// --- Public Auth routes (Epic 1)
	r.HandleFunc("/api/v1/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/login", authHandler.Login).Methods(http.MethodPost)

	// --- Public Product + Review routes (Epics 2 & 6)
	r.HandleFunc("/api/v1/products", productHandler.ListProducts).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/products/{id}", productHandler.GetProductByID).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/products/{id}/reviews", reviewHandler.ListReviews).Methods(http.MethodGet)

	//Protected routes
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware) // checks JWT, sets user in context

	// Cart routes (Epic 3)
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods(http.MethodGet)
	protected.HandleFunc("/cart", cartHandler.AddToCart).Methods(http.MethodPost)
	protected.HandleFunc("/cart/{product_id}", cartHandler.UpdateCartItem).Methods(http.MethodPatch)
	protected.HandleFunc("/cart/{product_id}", cartHandler.RemoveCartItem).Methods(http.MethodDelete)

	// Order routes (Epic 4)
	protected.HandleFunc("/checkout", orderHandler.Checkout).Methods(http.MethodPost)
	protected.HandleFunc("/orders", orderHandler.ListOrders).Methods(http.MethodGet)

	// Authenticated review submit/update (Epic 6.1)
	protected.HandleFunc("/products/{id}/reviews", reviewHandler.CreateOrUpdateReview).Methods(http.MethodPost)

	// Admin-only routes (Epic 5)
	admin := protected.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AdminMiddleware) // this will check role == "admin"

	admin.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	admin.HandleFunc("/products/{id}", productHandler.UpdateProduct).Methods(http.MethodPatch)

	return r
}
