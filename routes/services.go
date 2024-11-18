package routes

import (
	"PSbackend/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func ServiceRoutes(client *mongo.Client, dbName, serviceCollection, dbUserName, userCollection string, router *mux.Router) {
	// Define route for creating a new service for Back Office
	router.HandleFunc("/api/v1/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(client, dbName, serviceCollection, dbUserName, userCollection, w, r)
	}).Methods("POST")

	// Define route for creating a new service for Mobile App
	router.HandleFunc("/api/v1/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.CreateService(client, dbName, serviceCollection, dbUserName, userCollection, w, r)
	}).Methods("POST")

	// Define route for getting all services for Back Office
	router.HandleFunc("/api/v1/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all services for Mobile App
	router.HandleFunc("/api/v1/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.GetServices(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get a service by id for Back Office
	router.HandleFunc("/api/v1/bo/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get a service by id for Mobile App
	router.HandleFunc("/api/v1/mb/services/id", func(w http.ResponseWriter, r *http.Request) {
		api.GetService(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all filtered services by type for Back Office
	router.HandleFunc("/api/v1/bo/services/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetFilteredServiceType(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route for getting all filtered services by type for Mobile App
	router.HandleFunc("/api/v1/mb/services/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetFilteredServiceType(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to update a service for Back Office
	router.HandleFunc("/api/v1/bo/services", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateService(client, dbName, serviceCollection, w, r)
	}).Methods("PUT")

	// Define route to update a service for Mobile App
	router.HandleFunc("/api/v1/mb/services", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateService(client, dbName, serviceCollection, w, r)
	}).Methods("PUT")

	// Define route for creating a new specific service type for Back Office
	router.HandleFunc("/api/v1/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.CreateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route for getting all services types for Back Office
	router.HandleFunc("/api/v1/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route for getting all services types for Mobile App
	router.HandleFunc("/api/v1/mb/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route to update a service type for Back Office
	router.HandleFunc("/api/v1/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("PUT")

	// Define route to delete a service type by id for Back Office
	router.HandleFunc("/api/v1/bo/service-type", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteServiceType(client, dbName, os.Getenv("SERVICE_TYPE_COLLECTION"), w, r)
	}).Methods("DELETE")

	// Define route to get services by technician for mobile
	router.HandleFunc("/api/v1/mb/services/technicians", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceByTechnician(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to get services by technician for Back Office
	router.HandleFunc("/api/v1/bo/services/technicians", func(w http.ResponseWriter, r *http.Request) {
		api.GetServiceByTechnician(client, dbName, serviceCollection, w, r)
	}).Methods("GET")

	// Define route to update service with a new appointment for Mobile
	router.HandleFunc("/api/v1/mb/services/appointment", func(w http.ResponseWriter, r *http.Request) {
		api.InsertAppointment(client, dbName, serviceCollection, w, r)
	}).Methods("PUT")
}
