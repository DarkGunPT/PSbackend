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

// GetServices handles GET requests to get the list of services
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

// GetService handles GET requests to get one specific service
func GetService(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ID primitive.ObjectID `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var service models.Services
	collection := client.Database(dbName).Collection(serviceCollection)

	err = collection.FindOne(ctx, bson.M{"_id": requestBody.ID}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(service)
}

// DeleteService handles DELETE request to delete a specific service
func DeleteService(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		ID primitive.ObjectID `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(serviceCollection)

	result, err := collection.DeleteOne(ctx, bson.M{"_id": requestBody.ID})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Service deleted successfully")
}
