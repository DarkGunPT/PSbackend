package api

import (
	"PSbackend/models"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	if result.ModifiedCount == 0 && result.UpsertedID == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
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

// GetUser handles GET requests to get one specific user
func GetUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var requestBody struct {
		NIF int64 `json:"nif"`
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// UpdateUser handles PUT request to update one specific user
func UpdateUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)
	if user.Email == "" {
		http.Error(w, "Email is required for update", http.StatusBadRequest)
		return
	}

	collection := client.Database(dbName).Collection(userCollection)
	var updateFields bson.M = bson.M{}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	if user.NIF != 0 {
		updateFields["nif"] = user.NIF
	}
	if user.Phone != 0 {
		updateFields["phone"] = user.Phone
	}
	if len(user.Role) > 0 {
		updateFields["role"] = user.Role
	}
	if len(user.ServiceTypes) > 0 {
		updateFields["service_types"] = user.ServiceTypes
	}
	if user.Locality != "" {
		updateFields["locality"] = user.Locality
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": user.Email}, bson.M{"$set": updateFields}, options.Update().SetUpsert(true))
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
