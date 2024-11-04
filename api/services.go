package api

import (
	"PSbackend/models"
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// UpdateService handles PUT request to update one specific service
func UpdateService(ctx context.Context, client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var service models.Services

	json.NewDecoder(r.Body).Decode(&service)
	collection := client.Database(dbName).Collection(serviceCollection)
	var updateFields bson.M = bson.M{}
	if service.ServiceType != (models.ServiceType{}) {
		updateFields["service_type"] = service.ServiceType
	}
	if service.Description != "" {
		updateFields["description"] = service.Description
	}
	if len(service.Appointment) > 0 {
		updateFields["appointments"] = service.Appointment
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": service.ID}, bson.M{"$set": updateFields}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Service updated successfully")
}

// CreateServiceType handles POST requests to create a specific service type
func CreateServiceType(ctx context.Context, client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var serviceType models.ServiceType

	json.NewDecoder(r.Body).Decode(&serviceType)
	serviceType.ID = primitive.NewObjectID()

	collection := client.Database(dbName).Collection(serviceTypeCollection)

	result, err := collection.InsertOne(ctx, serviceType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetServiceType handles GET requests to get the list of service types
func GetServiceType(ctx context.Context, client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var servicesType []models.ServiceType

	collection := client.Database(dbName).Collection(serviceTypeCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var serviceType models.ServiceType
		cursor.Decode(&serviceType)
		servicesType = append(servicesType, serviceType)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servicesType)
}

// UpdateServiceType handles PUT request to update one specific service type
func UpdateServiceType(ctx context.Context, client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var serviceType models.ServiceType

	json.NewDecoder(r.Body).Decode(&serviceType)
	collection := client.Database(dbName).Collection(serviceTypeCollection)
	var updateFields bson.M = bson.M{}

	if serviceType.Name != "" {
		updateFields["name"] = serviceType.Name
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": serviceType.ID}, bson.M{"$set": updateFields}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Service updated successfully")
}
