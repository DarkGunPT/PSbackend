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

	// Define route to get a service by id
	router.HandleFunc("/api/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(ctx, client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to delete a service by id
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteService(ctx, client, dbName, serviceCollection, w, r)
	}).Methods("DELETE")

	// Define route to update a service
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateService(ctx, client, dbName, serviceCollection, w, r)
	}).Methods("PUT")
}
