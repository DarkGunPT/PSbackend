package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"PSbackend/config"
	"PSbackend/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRouter() (*http.ServeMux, *mongo.Client) {
	router := http.NewServeMux()
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	client, err := config.ConnectDB(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Error connecting to mongodb:", err)
	}

	// Retrieves all users or deletes a user
	router.HandleFunc("/api/v1/mb/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			GetUsers(client, "PSprojectTest", "Users", w, r)
		} else if r.Method == http.MethodDelete {
			DeleteUser(client, "PSprojectTest", "Users", w, r)
		}
	})

	// Retrieves all technicians
	router.HandleFunc("/api/v1/mb/users/technicians", func(w http.ResponseWriter, r *http.Request) {
		GetTechnicians(client, "PSprojectTest", "Users", w, r)
	})

	// Retrieves a user by NIF
	router.HandleFunc("/api/v1/bo/users/nif", func(w http.ResponseWriter, r *http.Request) {
		GetUser(client, "PSprojectTest", "Users", w, r)
	})

	// Retrieves all clients
	router.HandleFunc("/api/v1/mb/users/clients", func(w http.ResponseWriter, r *http.Request) {
		GetClients(client, "PSprojectTest", "Users", w, r)
	})

	// Retrieves clients ordered by a filter
	router.HandleFunc("/api/v1/bo/users/clients/order", func(w http.ResponseWriter, r *http.Request) {
		OrderClients(client, "PSprojectTest", "Users", w, r)
	})

	// Retrieves all fees
	router.HandleFunc("/api/v1/bo/fees", func(w http.ResponseWriter, r *http.Request) {
		GetFees(client, "PSprojectTest", "Users", w, r)
	})

	// Retrieves the count of services performed
	router.HandleFunc("/api/v1/bo/count-services-performed", func(w http.ResponseWriter, r *http.Request) {
		GetServicesPerformed(client, "PSprojectTest", "Users", w, r)
	})

	return router, client
}

func TestGetUsers(t *testing.T) {
	router, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/mb/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	router, client := setupRouter()

	newUser := models.User{
		ID:  primitive.NewObjectID(),
		NIF: 1234567890,
	}

	collection := client.Database("PSprojectTest").Collection("Users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, newUser)
	if err != nil {
		t.Fatalf("Failed to insert user: %v", err)
	}

	deleteData := map[string]interface{}{
		"nif": 1234567890,
	}
	deleteJSON, _ := json.Marshal(deleteData)

	// Test deleting an existing user
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/mb/users", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}

	// Test deleting a non-existent user
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/mb/users", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("Expected status code 404, got %v", w.Code)
	}
}

func TestGetClients(t *testing.T) {
	router, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/mb/users/clients", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}

func TestGetFees(t *testing.T) {
	router, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bo/fees", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}

func TestGetServicesPerformed(t *testing.T) {
	router, _ := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/bo/count-services-performed", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}
