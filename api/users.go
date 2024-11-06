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

	code := rand.Intn(8999) + 1000

	collection := mongo.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"email": requestBody.Email}, bson.M{"$set": bson.M{"recovery_code": code}}, options.Update().SetUpsert(true))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.ModifiedCount == 0 {
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

// CreateUser handles POST requests to create a new user
func CreateUser(client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)
	user.ID = primitive.NewObjectID()

	collection := client.Database(dbName).Collection(userCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	collection := client.Database(dbName).Collection(userCollection)
	var updateFields bson.M = bson.M{}
	if user.Name != "" {
		updateFields["name"] = user.Name
	}
	if user.Password != "" {
		updateFields["password"] = user.Password
	}
	if user.Phone != 0 {
		updateFields["phone"] = user.Phone
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
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
	if user.Rating != 0 {
		updateFields["rating"] = user.Rating
	}
	if user.BlockServices {
		updateFields["block_services"] = user.BlockServices
	}
	if user.IsActive {
		updateFields["is_active"] = user.IsActive
	}
	if !user.CreatedAt.IsZero() {
		updateFields["created_at"] = user.CreatedAt
	}
	if user.RecoveryCode != 0 {
		updateFields["recovery_code"] = user.RecoveryCode
	}

	if len(updateFields) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	result, err := collection.UpdateOne(ctx, bson.M{"nif": user.NIF}, bson.M{"$set": updateFields}, options.Update().SetUpsert(true))
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
	json.NewEncoder(w).Encode("Login successful")
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
