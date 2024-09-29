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

func main() {
	// Create context with timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
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
	routes.UserRoutes(ctx, client, os.Getenv("DB_NAME"), os.Getenv("USER_COLLECTION"), router)

	// Register service-related routes
	routes.ServiceRoutes(ctx, client, os.Getenv("DB_NAME"), os.Getenv("SERVICE_COLLECTION"), router)

	log.Println("Starting the http server at port :8080")
	// Start the HTTP server on port 8080
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting the http server:", err)
		return
	}
}
