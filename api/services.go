package api

import (
	"PSbackend/models"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetServices handles GET requests to get the list of services
func GetServices(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var user models.User
		cursor.Decode(&user)

		users = append(users, user)

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
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

	var service models.ServiceType
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
	var services []models.ServiceType

	var filter struct {
		ServiceType string `json:"name"`
	}

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	filter.ServiceType = strings.ToUpper(filter.ServiceType)

	collection := client.Database(dbName).Collection(serviceCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{"name": filter.ServiceType})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var service models.ServiceType
		cursor.Decode(&service)
		services = append(services, service)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

// UpdateService handles PUT request to update one specific service
func UpdateService(client *mongo.Client, dbName, serviceCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var service models.ServiceType

	json.NewDecoder(r.Body).Decode(&service)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(serviceCollection)
	var updateFields bson.M = bson.M{}

	if service.Price != 0 {
		updateFields["price"] = service.Price
	}

	if service.Name != "" {
		updateFields["name"] = strings.ToUpper(service.Name)
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
	serviceType.Name = strings.ToUpper(serviceType.Name)

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
		updateFields["name"] = strings.ToUpper(serviceType.Name)
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
	var services []models.ServiceType

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
		var service models.ServiceType
		cursor.Decode(&service)
		services = append(services, service)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(services)
}

func InsertAppointment(client *mongo.Client, dbName, serviceCollection, userCollection, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ClientEmail   string `json:"client_email"`
		ProviderEmail string `json:"provider_email"`
		ServiceName   string `json:"service_name"`
		Start         string `json:"start"`
		End           string `json:"end"`
		Phone         string `json:"phone"`
		NIF           string `json:"nif"`
		Locality      string `json:"locality"`
		Notes         string `json:"notes"`
		TotalPrice    string `json:"totalPrice"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	requestBody.ServiceName = strings.ToUpper(requestBody.ServiceName)

	var cli models.User

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.ClientEmail}).Decode(&cli)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var provider models.User

	collection = client.Database(dbName).Collection(userCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.ProviderEmail}).Decode(&provider)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
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

	phone, err := strconv.Atoi(requestBody.Phone)
	if err != nil {
		http.Error(w, "Invalid Phone format", http.StatusBadRequest)
		return
	}

	nif, err := strconv.Atoi(requestBody.NIF)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	totalPrice, err := strconv.ParseFloat(requestBody.TotalPrice, 64)
	if err != nil {
		http.Error(w, "Invalid Total price format", http.StatusBadRequest)
		return
	}

	var priceHour float64

	for _, service := range provider.ServiceTypes {
		if service.Name == strings.ToUpper(requestBody.ServiceName) {
			priceHour = service.Price
		}
	}

	appointment := models.Appointment{
		ID:          primitive.NewObjectID(),
		Provider:    provider,
		Client:      cli,
		Status:      "SCHEDULED",
		Start:       start,
		End:         end,
		Phone:       phone,
		NIF:         nif,
		Locality:    requestBody.Locality,
		Notes:       requestBody.Notes,
		TotalPrice:  totalPrice,
		PriceHour:   priceHour,
		ServiceName: requestBody.ServiceName,
	}

	for _, service := range provider.ServiceTypes {
		if service.Name == requestBody.ServiceName {
			appointment.PriceHour = service.Price
		}
	}

	collection = client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, appointment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var updateCli bson.M = bson.M{}
	for i := range cli.Role {
		if cli.Role[i].Name == "CLIENT" {
			cli.Role[i].ServicesDone++
		}
	}

	updateCli["role"] = cli.Role
	var updateTech bson.M = bson.M{}
	for i := range provider.Role {
		if provider.Role[i].Name == "TECH" {
			provider.Role[i].ServicesDone++
		}
	}

	updateTech["role"] = provider.Role
	collection = client.Database(dbName).Collection(userCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	update, err := collection.UpdateOne(ctx, bson.M{"email": cli.Email}, bson.M{"$set": updateCli})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if update.MatchedCount == 0 {
		http.Error(w, "Client not found to update", http.StatusNotFound)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	update, err = collection.UpdateOne(ctx, bson.M{"email": provider.Email}, bson.M{"$set": bson.M{"role": provider.Role}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if update.MatchedCount == 0 {
		http.Error(w, "Provider not found to update", http.StatusNotFound)
		return
	}

	jsonResponse := map[string]interface{}{
		"message": "Appointment created successfully",
		"result":  result,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

// GetAppointments handles GET requests to get the list of appointments
func GetAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var appointments []models.Appointment

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)
		appointments = append(appointments, appointment)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetUpcommingAppointments handles GET requests to get the list of appointments
func GetUpcommingAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var appointments []models.Appointment

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Status == "SCHEDULED" {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetClientUpcommingAppointments handles GET requests to get the list of upcomming appointments of a client
func GetClientUpcommingAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nif, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}
	nifInt, err := strconv.Atoi(nif)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Client.NIF == nifInt && appointment.Status == "SCHEDULED" {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetTechUpcommingAppointments handles GET requests to get the list of upcomming appointments of a tech
func GetTechUpcommingAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nif, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}
	nifInt, err := strconv.Atoi(nif)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Provider.NIF == nifInt && appointment.Status == "SCHEDULED" {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetClientHistoryAppointments handles GET requests to get the list of appointments of a client already CLOSED
func GetClientHistoryAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nif, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}
	nifInt, err := strconv.Atoi(nif)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Client.NIF == nifInt && appointment.Status == "COMPLETED" || appointment.Status == "CANCELED" {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetTechHistoryAppointments handles GET requests to get the list of appointments of a tech already closed
func GetTechHistoryAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nif, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}
	nifInt, err := strconv.Atoi(nif)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Provider.NIF == nifInt && appointment.Status == "COMPLETED" || appointment.Status == "CANCELED" {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetHistoryAppointments handles GET requests to get the list of appointments already CLOSED
func GetHistoryAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if time.Now().After(appointment.End) {
			appointment.Status = "COMPLETED"
			filter := bson.M{"_id": appointment.ID}
			update := bson.M{"$set": bson.M{"status": "COMPLETED"}}
			_, err := collection.UpdateOne(ctx, filter, update)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		if appointment.Status == "COMPLETED" || appointment.Status == "CANCELED" {
			appointments = append(appointments, appointment)
		}

	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetHistoryAppointments handles GET requests to get the list of appointments already CLOSED
func GetAppointmentsByPrice(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ServiceType string  `json:"service_type"`
		Max         float64 `json:"max"`
		Min         float64 `json:"min"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	appointments := make([]models.Appointment, 0)

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)

		if appointment.ServiceName == requestBody.ServiceType && appointment.TotalPrice >= requestBody.Min && appointment.TotalPrice <= requestBody.Max {
			appointments = append(appointments, appointment)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(appointments)
}

// GetServicesByPrice handles GET requests to get the list of appointments already CLOSED
func GetServicesByPrice(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		ServiceType string  `json:"service_type"`
		Max         float64 `json:"max"`
		Min         float64 `json:"min"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var temporaryUser models.User
	resultUsers := make([]models.User, 0)
	for cursor.Next(ctx) {
		var user models.User
		cursor.Decode(&user)
		temporaryUser = user

		for _, service := range user.ServiceTypes {
			if service.Name == requestBody.ServiceType && service.Price >= requestBody.Min && service.Price <= requestBody.Max {
				temporaryUser.ServiceTypes = []models.ServiceType{}
				temporaryUser.ServiceTypes = append(temporaryUser.ServiceTypes, service)
				resultUsers = append(resultUsers, temporaryUser)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultUsers)
}

// GetServicesByPrice handles GET requests to get the list of appointments already CLOSED
func GetServicesByPriceQuery(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	serviceType := r.URL.Query().Get("service_type")
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	min, err := strconv.ParseFloat(minStr, 64)
	if err != nil {
		http.Error(w, "Invalid min value", http.StatusBadRequest)
		return
	}

	max, err := strconv.ParseFloat(maxStr, 64)
	if err != nil {
		http.Error(w, "Invalid max value", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var temporaryUser models.User
	resultUsers := make([]models.User, 0)
	for cursor.Next(ctx) {
		var user models.User
		cursor.Decode(&user)
		temporaryUser = user

		for _, service := range user.ServiceTypes {
			if service.Name == serviceType && service.Price >= min && service.Price <= max {
				temporaryUser.ServiceTypes = []models.ServiceType{}
				temporaryUser.ServiceTypes = append(temporaryUser.ServiceTypes, service)
				resultUsers = append(resultUsers, temporaryUser)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resultUsers)
}

func DeleteAppointment(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, exists := vars["id"]
	if !exists {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"status": "CANCELED"}}
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Appointment not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Appointment canceled successfully")
}

// GetAppointments handles GET requests to get the list of appointments
func GetCountAppointments(client *mongo.Client, dbName, appointmentCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(appointmentCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var total int

	for cursor.Next(ctx) {
		var appointment models.Appointment
		cursor.Decode(&appointment)
		total++
	}

	jsonResponse := map[string]interface{}{
		"message": "Counted every existent Appointment",
		"count":   total,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}
