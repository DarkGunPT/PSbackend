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
)

// setupRouter configura o roteador para testes
func setupRouter() *http.ServeMux {
	router := http.NewServeMux()
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	client, err := config.ConnectDB(ctx, "mongodb://localhost:27017")
	if err != nil {
		log.Fatal("Error connecting to mongodb:", err)
	}

	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateUser(context.Background(), client, "PSproject", "Users", w, r)
		} else if r.Method == http.MethodGet {
			GetUsers(context.Background(), client, "PSproject", "Users", w, r)
		} else if r.Method == http.MethodDelete {
			DeleteUser(context.Background(), client, "PSproject", "Users", w, r)
		}
	})

	router.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		Login(context.Background(), client, "PSproject", "Users", w, r)
	})

	router.HandleFunc("/api/users/phone", func(w http.ResponseWriter, r *http.Request) {
		GetUser(context.Background(), client, "PSproject", "Users", w, r)
	})

	return router
}

func TestCreateUser(t *testing.T) {
	router := setupRouter()

	user := models.User{
		ID:           primitive.NewObjectID(),
		Phone:        1234567890,
		Name:         "John Doe",
		Password:     "password123",
		Role:         "freelancer",
		Locality:     "Downtown",
		Profession:   "Developer",
		PricePerHour: 50.0,
		Rating:       4.5,
		Fee:          10.0,
		JobsDone:     5,
		Block:        false,
	}
	userJSON, _ := json.Marshal(user)

	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}

func TestGetUsers(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}

func TestLoginUser(t *testing.T) {
	router := setupRouter()

	loginData := map[string]interface{}{
		"phone":    1234567890,
		"password": "password123",
	}
	loginJSON, _ := json.Marshal(loginData)

	req := httptest.NewRequest(http.MethodPost, "/api/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v. Response: %s", w.Code, w.Body.String())
	}
}

func TestDeleteUser(t *testing.T) {
	router := setupRouter()

	deleteData := map[string]interface{}{
		"phone": 1234567890,
	}
	deleteJSON, _ := json.Marshal(deleteData)

	req := httptest.NewRequest(http.MethodDelete, "/api/users", bytes.NewBuffer(deleteJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status code 200, got %v", w.Code)
	}
}
