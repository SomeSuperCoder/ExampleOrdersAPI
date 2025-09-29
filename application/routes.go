package application

import (
	"fmt"
	"net/http"

	"github.com/SomeSuperCoder/OrdersAPI/handlers"
	"github.com/SomeSuperCoder/OrdersAPI/middleware"
	"github.com/SomeSuperCoder/OrdersAPI/repository"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func loadRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Ok")
	})
	mux.Handle("/orders/", loadOrderRoutes(db))
	mux.Handle("/orders", loadOrderRoutes(db))

	return middleware.LoggerMiddleware(mux)
}

func loadOrderRoutes(db *mongo.Database) http.Handler {
	mux := http.NewServeMux()
	orderHandler := handlers.OrderHandler{
		Repo: repository.NewOrderRepo(db),
	}

	mux.HandleFunc("GET /{id}", orderHandler.GetOrder)
	mux.HandleFunc("POST /", orderHandler.CreateOrder)
	mux.HandleFunc("PATCH /{id}", orderHandler.UpdateOrder)

	return http.StripPrefix("/orders", mux)
}
