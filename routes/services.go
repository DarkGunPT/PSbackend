package routes

import (
	"PSbackend/api"
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceRoutes(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, router *mux.Router) {
	// Define route for creating a new service
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(ctx, client, dbName, serviceCollection, w, r)
	}).Methods("POST")

	// Define route for getting all services
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(ctx, client, dbName, serviceCollection, w, r)
	}).Methods("GET")
}
