package api

import (
	"PSbackend/models"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/signintech/gopdf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func emailExists(email string, client *mongo.Client, dbName, userCollection string) bool {
	// Check if the email exists in the database
	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var existingUser struct {
		Email string `bson:"email"`
	}
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&existingUser)
	return err == nil
}

// RecoveryEmail is responsible for sending an recovery email with a verification code to the user's email address
func RecoveryEmail(w http.ResponseWriter, r *http.Request, mongo *mongo.Client, dbName, userCollection string) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !emailExists(requestBody.Email, mongo, dbName, userCollection) {
		http.Error(w, "Email isn't registered", http.StatusConflict)
		return
	}

	code := rand.Intn(8999) + 1000

	collection := mongo.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": requestBody.Email}, bson.M{"$set": bson.M{"recovery_code": code}}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 && result.UpsertedID == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	from := mail.NewEmail("FixFinder", os.Getenv("FIXFINDER_EMAIL"))
	to := mail.NewEmail("Nuno Honório", requestBody.Email)
	subject := fmt.Sprintf("Email Validation of FixFinder Code: %d", code)
	plainTextContent := "Making it easier to find technicians for certain domestic services​"
	htmlContent := "<strong>Making it easier to find technicians for certain domestic services​</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_APIKEY"))
	response, err := client.Send(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response.StatusCode >= 400 {
		http.Error(w, fmt.Sprintf("Failed to send email, status code: %d", response.StatusCode), http.StatusInternalServerError)
		return
	}

	// Send the generated code back to the frontend in the response
	jsonResponse := map[string]interface{}{
		"message": "Verification email sent successfully",
		"code":    code,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

// VerificateEmail is responsible for sending an email with a verification code to the user's email address
func VerificateEmail(w http.ResponseWriter, r *http.Request, mongo *mongo.Client, dbName, userCollection string) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if emailExists(requestBody.Email, mongo, dbName, userCollection) {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	code := rand.Intn(8999) + 1000

	defaultUser := models.User{
		Name:          "", // Default name
		Password:      "", // Default password
		NIF:           0,
		Phone:         0,
		Email:         requestBody.Email, // Use email from the request
		Role:          []models.Role{},
		ServiceTypes:  []models.ServiceType{},
		Locality:      "",
		Rating:        0.0,
		BlockServices: false,
		IsActive:      false,
		CreatedAt:     time.Now(), // Default to current timestamp
		RecoveryCode:  code,
	}

	collection := mongo.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": requestBody.Email}, bson.M{"$set": defaultUser}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 && result.UpsertedID == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	from := mail.NewEmail("FixFinder", os.Getenv("FIXFINDER_EMAIL"))
	to := mail.NewEmail("Nuno Honório", requestBody.Email)
	subject := fmt.Sprintf("Email Validation of FixFinder Code: %d", code)
	plainTextContent := "Making it easier to find technicians for certain domestic services​"
	htmlContent := "<strong>Making it easier to find technicians for certain domestic services​</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_APIKEY"))
	response, err := client.Send(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if response.StatusCode >= 400 {
		http.Error(w, fmt.Sprintf("Failed to send email, status code: %d", response.StatusCode), http.StatusInternalServerError)
		return
	}

	// Send the generated code back to the frontend in the response
	jsonResponse := map[string]interface{}{
		"message": "Verification email sent successfully",
		"code":    code,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

// Confirm the given code from user with the one saved in database
func ConfirmAuthCode(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email    string `json:"email"`
		Code     int    `json:"code"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user.RecoveryCode != requestBody.Code {
		http.Error(w, "Incorrect email verification code", http.StatusUnauthorized)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": requestBody.Email}, bson.M{"$set": bson.M{"password": requestBody.Password}}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Password reseted successfully")
}

// GetUsers handles GET requests to get the list of users
func GetUsers(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
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

// GetUser handles GET requests to get one specific user by NIF
func GetUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nifStr, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}

	nif, err := strconv.ParseInt(nifStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	var user models.User
	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"nif": nif}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT request to update one specific user
func UpdateUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	nifStr, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}

	nif, err := strconv.ParseInt(nifStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Name         string               `json:"name" bson:"name"`
		Password     string               `json:"password" bson:"password"`
		Phone        string               `json:"phone" bson:"phone"`
		Role         []models.Role        `json:"role" bson:"role"`
		ServiceTypes []models.ServiceType `json:"service_types" bson:"service_types"`
		Locality     string               `json:"locality" bson:"locality"`
		WorkStart    string               `json:"workStart" bson:"workStart"`
		WorkEnd      string               `json:"workEnd" bson:"workEnd"`
	}

	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	parsedPhone, err := strconv.ParseInt(requestBody.Phone, 10, 64)
	if err != nil {
		http.Error(w, "Invalid Phone format", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)
	var updateFields bson.M = bson.M{}
	if requestBody.Name != "" {
		updateFields["name"] = requestBody.Name
	}
	if requestBody.Password != "" {
		updateFields["password"] = requestBody.Password
	}
	if parsedPhone != 0 {
		updateFields["phone"] = parsedPhone
	}
	if len(requestBody.Role) > 0 {
		updateFields["role"] = requestBody.Role
	}
	if len(requestBody.ServiceTypes) > 0 {
		updateFields["service_types"] = requestBody.ServiceTypes
	}
	if requestBody.Locality != "" {
		updateFields["locality"] = requestBody.Locality
	}
	if requestBody.WorkStart != "" {
		start, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.WorkStart)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updateFields["workStart"] = start
	}
	if requestBody.WorkEnd != "" {
		end, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.WorkEnd)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		updateFields["workEnd"] = end
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"nif": nif}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"nif": nif}, bson.M{"$set": updateFields}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "User wasn't modified", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User updated successfully")
}

// DeleteUser handles DELETE request to delete a specific user
func DeleteUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		NIF int64 `json:"nif"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(userCollection)

	result, err := collection.DeleteOne(ctx, bson.M{"nif": requestBody.NIF})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User deleted successfully")
}

func Login(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if user.Password != requestBody.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func LoginAdmin(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	var isAdmin bool
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)

	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for _, role := range user.Role {
		if role.Name == "ADMIN" {
			isAdmin = true
		}
	}

	if !isAdmin {
		http.Error(w, "User isn't admin", http.StatusUnauthorized)
		return
	}

	if user.Password != requestBody.Password {
		http.Error(w, "Incorrect password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Login successful")
}

// CreateRole handles POST requests to create a new role
func CreateRole(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var role models.Role

	json.NewDecoder(r.Body).Decode(&role)
	role.ID = primitive.NewObjectID()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(userCollection)

	result, err := collection.InsertOne(ctx, role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// GetTechnicians handles GET requests to get the list of technicians
func GetTechnicians(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User

		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "TECH" {
				users = append(users, user)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetClients handles GET requests to get the list of clients
func GetClients(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User

		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "CLIENT" {
				users = append(users, user)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// UpdateUser handles PUT request to update one specific user
func RegisterCompletion(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var requestBody struct {
		Email        string               `json:"email" bson:"email"`
		Name         string               `json:"name" bson:"name"`
		NIF          string               `json:"nif" bson:"nif"`
		Phone        string               `json:"phone" bson:"phone"`
		ServiceTypes []models.ServiceType `json:"service_types" bson:"service_types"`
		Locality     string               `json:"locality" bson:"locality"`
		WorkStart    string               `json:"workStart" bson:"workStart"`
		WorkEnd      string               `json:"workEnd" bson:"workEnd"`
	}

	json.NewDecoder(r.Body).Decode(&requestBody)
	if requestBody.Email == "" {
		http.Error(w, "Email is required for update", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)
	if requestBody.Name == "" {
		http.Error(w, "Name is required for registration", http.StatusBadRequest)
		return
	}
	if requestBody.NIF == "" {
		http.Error(w, "NIF is required for registration", http.StatusBadRequest)
		return
	}
	if requestBody.Phone == "" {
		http.Error(w, "Phone is required for registration", http.StatusBadRequest)
		return
	}
	if requestBody.Locality == "" {
		http.Error(w, "Locality is required for registration", http.StatusBadRequest)
		return
	}

	nif, err := strconv.Atoi(requestBody.NIF)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	phone, err := strconv.Atoi(requestBody.Phone)
	if err != nil {
		http.Error(w, "Invalid Phone format", http.StatusBadRequest)
		return
	}

	var updateFields bson.M = bson.M{}
	updateFields["name"] = requestBody.Name
	updateFields["nif"] = nif
	updateFields["phone"] = phone
	if len(requestBody.ServiceTypes) == 0 {
		updateFields["role"] = []models.Role{
			{
				Name: "CLIENT",
			},
		}
	} else {
		updateFields["role"] = []models.Role{
			{
				Name: "CLIENT",
			},
			{
				Name: "TECH",
			},
		}
	}

	start, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.WorkStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updateFields["workStart"] = start
	end, err := time.Parse("2006-01-02T15:04:05.999-07:00", requestBody.WorkEnd)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	updateFields["workEnd"] = end
	updateFields["service_types"] = requestBody.ServiceTypes
	updateFields["locality"] = requestBody.Locality
	updateFields["is_active"] = true

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var updatedUser models.User
	err = collection.FindOneAndUpdate(ctx, bson.M{"email": requestBody.Email}, bson.M{"$set": updateFields}, opts).Decode(&updatedUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func UpdateBlock(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !user.BlockServices {
		user.BlockServices = true
	} else {
		user.BlockServices = false
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": user}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User updated successfully")

}

func UpdateActive(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Email string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, bson.M{"email": requestBody.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !user.IsActive {
		user.IsActive = true
	} else {
		user.IsActive = false
	}

	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": user}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("User updated successfully")

}

// GetTechnicians handles GET requests to get the list of technicians
func OrderTechnicians(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Filter string `json:"filter"`
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

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User

		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "TECH" {
				users = append(users, user)
			}
		}
	}

	if requestBody.Filter == "rating" {
		sort.Slice(users, func(i, j int) bool {
			return users[i].Rating > users[j].Rating
		})
	}
	if requestBody.Filter == "services" {
		sort.Slice(users, func(i, j int) bool {
			var servicesDoneI, servicesDoneJ int

			for _, role := range users[i].Role {
				if role.Name == "TECH" {
					servicesDoneI = role.ServicesDone
					break
				}
			}

			for _, role := range users[j].Role {
				if role.Name == "TECH" {
					servicesDoneJ = role.ServicesDone
					break
				}
			}

			return servicesDoneI > servicesDoneJ
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetClients handles GET requests to get the list of clients
func OrderClients(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		Filter string `json:"filter"`
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

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User

		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "CLIENT" {
				users = append(users, user)
			}
		}
	}

	if requestBody.Filter == "rating" {
		sort.Slice(users, func(i, j int) bool {
			return users[i].Rating > users[j].Rating
		})
	}
	if requestBody.Filter == "services" {
		sort.Slice(users, func(i, j int) bool {
			var servicesDoneI, servicesDoneJ int

			for _, role := range users[i].Role {
				if role.Name == "CLIENT" {
					servicesDoneI = role.ServicesDone
					break
				}
			}

			for _, role := range users[j].Role {
				if role.Name == "CLIENT" {
					servicesDoneJ = role.ServicesDone
					break
				}
			}

			return servicesDoneI > servicesDoneJ
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}

// GetUser handles GET requests to get one specific user by NIF
func GetFeesByNif(client *mongo.Client, dbName, feesCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)

	nifStr, exists := vars["nif"]
	if !exists {
		http.Error(w, "NIF is required", http.StatusBadRequest)
		return
	}

	nif, err := strconv.ParseInt(nifStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid NIF format", http.StatusBadRequest)
		return
	}

	fees := make([]models.Fee, 0)

	collection := client.Database(dbName).Collection(feesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var fee models.Fee
		cursor.Decode(&fee)
		if fee.NIF == nif {
			fees = append(fees, fee)
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fees)
}

// GetUser handles GET requests to get one specific user by NIF
func GetFees(client *mongo.Client, dbName, feesCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(feesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	fees := make([]models.Fee, 0)

	for cursor.Next(ctx) {
		var fee models.Fee
		cursor.Decode(&fee)
		fees = append(fees, fee)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fees)
}

func CreateFee(client *mongo.Client, dbName, feesCollection string, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		NIF   int     `json:"nif"`
		Value float64 `json:"value"`
		Day   string  `json:"day"`
		Month string  `json:"month"`
		Year  string  `json:"year"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"nif": requestBody.NIF}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var servicesDone int

	for _, role := range user.Role {
		if role.Name == "TECH" {
			servicesDone = role.ServicesDone
		}
	}

	fee := models.Fee{
		NIF:      int64(requestBody.NIF),
		Value:    requestBody.Value,
		JobsDone: int64(servicesDone),
		Paid:     false,
		Day:      requestBody.Day,
		Month:    requestBody.Month,
		Year:     requestBody.Year,
	}

	collection = client.Database(dbName).Collection(feesCollection)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, fee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func PayFee(client *mongo.Client, dbName, feesCollection, userCollection string, w http.ResponseWriter, r *http.Request) {
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

	collection := client.Database(dbName).Collection(feesCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	update := bson.M{"$set": bson.M{"paid": true}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After).SetUpsert(true)
	var updatedFee models.Fee
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": objectID}, update, opts).Decode(&updatedFee)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	generateInvoice(updatedFee, client, dbName, userCollection)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Fee paid successfully")
}

func generateInvoice(invoice models.Fee, client *mongo.Client, dbName, userCollection string) {
	var user models.User
	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.M{"nif": invoice.NIF}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Print(err.Error())
			return
		}
		log.Print(err.Error())
		return
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	fontPath := filepath.Join(currentDir, "api", "fonts", "BebasNeue-Regular.ttf")
	err = pdf.AddTTFFont("BebasNeue-Regular", fontPath)
	if err != nil {
		log.Print(err.Error())
		return
	}

	err = pdf.SetFont("BebasNeue-Regular", "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}

	pdf.Cell(nil, "FixFinder - Invoice of Monthly Fee")
	pdf.Br(20)
	pdf.Cell(nil, fmt.Sprintf("%s of %s, %s", invoice.Day, invoice.Month, invoice.Year))
	pdf.Br(20)
	pdf.Cell(nil, fmt.Sprintf("NIF: %d", invoice.NIF))
	pdf.Br(20)
	pdf.Cell(nil, fmt.Sprintf("Value: %f", invoice.Value))
	pdf.Br(20)
	pdf.Cell(nil, "Status: Paid")
	pdfName := fmt.Sprintf("%s.pdf", invoice.ID.Hex())
	pdfPath := filepath.Join(currentDir, "api", "invoices", pdfName)
	pdf.WritePdf(pdfPath)

	file, err := os.Open(pdfPath) // Replace with your invoice path
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileContent := make([]byte, fileInfo.Size())
	_, err = file.Read(fileContent)
	if err != nil {
		log.Fatal(err)
	}

	encodedFile := base64.StdEncoding.EncodeToString(fileContent)

	attachment := mail.NewAttachment()
	attachment.SetContent(encodedFile)
	attachment.SetType("application/pdf")
	attachment.SetFilename(pdfName)
	attachment.SetDisposition("attachment")
	from := mail.NewEmail("FixFinder", os.Getenv("FIXFINDER_EMAIL"))
	to := mail.NewEmail("Provider", user.Email)
	subject := fmt.Sprintf("FixFinder Invoice Fee %s, %s", invoice.Month, invoice.Year)
	plainTextContent := "Making it easier to find technicians for certain domestic services​"
	htmlContent := "<strong>Making it easier to find technicians for certain domestic services​</strong>"
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	message.AddAttachment(attachment)
	sendgrid := sendgrid.NewSendClient(os.Getenv("SENDGRID_APIKEY"))
	response, err := sendgrid.Send(message)
	if err != nil {
		log.Print(err.Error())
		return
	}

	if response.StatusCode >= 400 {
		log.Printf("Failed to send email, status code: %d", response.StatusCode)
		return
	}
}

// GetUsers handles GET requests to get the list of users
func GetServicesPerformed(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(userCollection)
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
		var user models.User
		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "TECH" {
				total = total + role.ServicesDone
			}
		}
	}

	jsonResponse := map[string]interface{}{
		"message": "Counted every Service performed",
		"count":   total,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}

// GetUsers handles GET requests to get the list of users
func GetServicesReceived(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	collection := client.Database(dbName).Collection(userCollection)
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
		var user models.User
		cursor.Decode(&user)
		for _, role := range user.Role {
			if role.Name == "CLIENT" {
				total = total + role.ServicesDone
			}
		}
	}

	jsonResponse := map[string]interface{}{
		"message": "Counted every Service received",
		"count":   total,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(jsonResponse)
}
