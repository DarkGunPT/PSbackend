package main

import (
	"PSbackend/config"
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Create context with timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

	log.Println("Connected and Tested MongoDB")
}
