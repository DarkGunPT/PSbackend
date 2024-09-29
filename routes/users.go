package routes

import (
	"PSbackend/api"
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(ctx context.Context, client *mongo.Client, dbName, userCollection string, router *mux.Router) {
	// Define route for creating a new user
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.CreateUser(ctx, client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to get users
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.GetUsers(ctx, client, dbName, userCollection, w, r)
	}).Methods("GET")
}
