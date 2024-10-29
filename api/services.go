package api

import (
	"PSbackend/models"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateService handles POST requests to create a new service
func CreateService(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var service models.Services

	json.NewDecoder(r.Body).Decode(&service)
	service.ID = primitive.NewObjectID()

	collection := client.Database(dbName).Collection(serviceCollection)

	result, err := collection.InsertOne(ctx, service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

// GetServices handles GET requests to get the list of users
func GetServices(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var services []models.Services

	collection := client.Database(dbName).Collection(serviceCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var service models.Services
		cursor.Decode(&service)
		services = append(services, service)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}
