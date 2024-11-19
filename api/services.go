package api

import (
	"PSbackend/models"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateService handles POST requests to create a new service
func CreateService(client *mongo.Client, dbServiceName, serviceCollection, dbUserName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ServiceType models.ServiceType `json:"service_type"`
		Description string             `json:"description"`
		Employee_id string             `json:"employee_id"`
	}
	var service models.Services
	var technician models.User

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbUserName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Employee_id}).Decode(&technician)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.ID = primitive.NewObjectID()
	service.Employee = technician
	service.Description = requestBody.Description
	service.ServiceType = requestBody.ServiceType

	collection = client.Database(dbServiceName).Collection(serviceCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

// GetServices handles GET requests to get the list of services
func GetServices(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var services []models.Services

	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
func GetService(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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

// GetServiceType handles GET requests to get the filtered list of services by type
func GetFilteredServiceType(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var services []models.Services

	var filter struct {
		ServiceType string `json:"service_type"`
	}

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"service_type.name": filter.ServiceType})

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

// UpdateService handles PUT request to update one specific service
func UpdateService(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var service models.Services

	json.NewDecoder(r.Body).Decode(&service)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
func CreateServiceType(client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var serviceType models.ServiceType

	json.NewDecoder(r.Body).Decode(&serviceType)
	serviceType.ID = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
func GetServiceType(client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var servicesType []models.ServiceType

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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
func UpdateServiceType(client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var serviceType models.ServiceType

	json.NewDecoder(r.Body).Decode(&serviceType)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
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

// DeleteServiceType handles DELETE request to delete a specific service type
func DeleteServiceType(client *mongo.Client, dbName, serviceTypeCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		ID primitive.ObjectID `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(serviceTypeCollection)

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

// GetServiceByTechnician handles GET requests to get one specific service
func GetServiceByTechnician(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var services []models.Services

	var filter struct {
		EmployeeID string `json:"employee_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{
		"employee_id": filter.EmployeeID,
	})

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

func InsertAppointment(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ID         primitive.ObjectID `json:"id"`
		EmployerID string             `json:"employer_id"`
		Start      string             `json:"start"`
		End        string             `json:"end"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var service models.Services
	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"_id": requestBody.ID}).Decode(&service)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Service not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.Start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	end, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.End)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	service.Appointment = append(service.Appointment, models.Appointment{
		ID:         primitive.NewObjectID(),
		EmployeeID: service.Employee.ID.Hex(),
		EmployerID: requestBody.EmployerID,
		Status:     "CREATED",
		Start:      start,
		End:        end,
	})

	result, err := collection.UpdateOne(ctx, bson.M{"_id": service.ID}, bson.M{"$set": service}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Service appointments updated successfully")
}

// GetAppointments handles GET requests to get the list of appointments
func GetAppointments(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var appointments []models.Appointment

	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var service models.Services
		cursor.Decode(&service)
		appointments = append(appointments, service.Appointment...)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}
