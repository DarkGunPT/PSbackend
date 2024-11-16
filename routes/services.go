package routes

import (
	"PSbackend/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceRoutes(client *mongo.Client, dbName, serviceCollection string, router *mux.Router) {
	// Define route for creating a new service for Back Office
	router.HandleFunc("/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(client, dbName, serviceCollection, w, r)
	}).Methods("POST")

	// Define route for creating a new service for Mobile App
	router.HandleFunc("/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(client, dbName, serviceCollection, w, r)
	}).Methods("POST")

	// Define route for getting all services for Back Office
	router.HandleFunc("/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all services for Mobile App
	router.HandleFunc("/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get a service by id for Back Office
	router.HandleFunc("/bo/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get a service by id for Mobile App
	router.HandleFunc("/mb/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all filtered services by type for Back Office
	router.HandleFunc("/bo/services/type", func(w http.ResponseWriter, r *http.Request) {
		api.GetFilteredServiceType(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all filtered services by type for Mobile App
	router.HandleFunc("/mb/services/type", func(w http.ResponseWriter, r *http.Request) {
		api.GetFilteredServiceType(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to delete a service by id for Back Office
	router.HandleFunc("/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteService(client, dbName, serviceCollection, w, r)
	}).Methods("DELETE")

	// Define route to delete a service by id for Mobile App
	router.HandleFunc("/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteService(client, dbName, serviceCollection, w, r)
	}).Methods("DELETE")

	// Define route to update a service
	router.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateService(client, dbName, serviceCollection, w, r)
	}).Methods("PUT")

	// Define route for creating a new specific service type for Back Office
	router.HandleFunc("/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.CreateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route for creating a new specific service type for Mobile App
	router.HandleFunc("/mb/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.CreateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route for getting all services types for Back Office
	router.HandleFunc("/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route for getting all services types for Mobile App
	router.HandleFunc("/mb/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route to update a service type for Back Office
	router.HandleFunc("/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("PUT")

	// Define route to update a service type for Mobile App
	router.HandleFunc("/mb/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("PUT")

	// Define route to delete a service type by id for Back Office
	router.HandleFunc("/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("DELETE")

	// Define route to delete a service type by id for Mobile App
	router.HandleFunc("/mb/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("DELETE")
}
