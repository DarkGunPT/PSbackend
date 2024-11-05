package routes

import (
	"PSbackend/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceRoutes(client *mongo.Client, dbName, serviceCollection string, router *mux.Router) {
	// Define route for creating a new service
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(client, dbName, serviceCollection, w, r)
	}).Methods("POST")

	// Define route for getting all services
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get a service by id
	router.HandleFunc("/api/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to delete a service by id
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteService(client, dbName, serviceCollection, w, r)
	}).Methods("DELETE")

	// Define route to update a service
	router.HandleFunc("/api/services", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateService(client, dbName, serviceCollection, w, r)
	}).Methods("PUT")

	// Define route for creating a new specific service type
	router.HandleFunc("/api/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route for getting all services types
	router.HandleFunc("/services/type", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route to update a service type
	router.HandleFunc("/services/type", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("PUT")

	// Define route to delete a service type by id
	router.HandleFunc("/services/type", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("DELETE")
}
