package main

import (
	"PSbackend/config"
	"PSbackend/routes"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// disableCORS Middleware function to disable CORS
func disableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create context with timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	// Initialize the MongoDB connection
	client, err := config.ConnectDB(ctx, os.Getenv("MONGO_URI"))
	if err != nil {
		log.Fatal("Error connecting to mongodb:", err)
		return
	}

	err = config.TestConnection(ctx, *client)
	if err != nil {
		log.Fatal("Error testing connection with mongodb:", err)
		return
	}

	log.Println("Connected and Tested MongoDB with success")

	// Initialize Gorilla Mux router
	router := mux.NewRouter()

	// Register user-related routes
	routes.UserRoutes(client, os.Getenv("DB_NAME"), os.Getenv("USER_COLLECTION"), router)

	// Register service-related routes
	routes.ServiceRoutes(client, os.Getenv("DB_NAME"), os.Getenv("SERVICE_COLLECTION"), router)

	log.Println("Starting the http server at port :8080")
	// Start the HTTP server on port 8080
	err = http.ListenAndServe(":8080", disableCORS(router))
	if err != nil {
		log.Fatal("Error starting the http server:", err)
		return
	}
}
