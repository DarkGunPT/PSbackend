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

// CreateUser handles POST requests to create a new user
func CreateUser(ctx context.Context, client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User

	json.NewDecoder(r.Body).Decode(&user)
	user.ID = primitive.NewObjectID()

	collection := client.Database(dbName).Collection(userCollection)

	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

// GetUsers handles GET requests to get the list of users
func GetUsers(ctx context.Context, client *mongo.Client, dbName, userCollection string, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []models.User

	collection := client.Database(dbName).Collection(userCollection)

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
	json.NewEncoder(w).Encode(users)
}
