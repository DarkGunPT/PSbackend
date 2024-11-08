package routes

import (
	"PSbackend/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(client *mongo.Client, dbName, userCollection string, router *mux.Router) {
	// Define route for creating a new user
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.CreateUser(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to get users
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.GetUsers(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get a user by phone
	router.HandleFunc("/api/users/phone", func(w http.ResponseWriter, r *http.Request) {
		api.GetUser(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to update a user
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateUser(client, dbName, userCollection, w, r)
	}).Methods("PUT")

	// Define route to delete a user by phone
	router.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteUser(client, dbName, userCollection, w, r)
	}).Methods("DELETE")

	// Define route for user login
	router.HandleFunc("/api/users/login", func(w http.ResponseWriter, r *http.Request) {
		api.Login(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route email verification of mb
	router.HandleFunc("/api/mb/users/recovery", func(w http.ResponseWriter, r *http.Request) {
		api.VerificateEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route email verification of bo
	router.HandleFunc("/api/bo/users/verification", func(w http.ResponseWriter, r *http.Request) {
		api.VerificateEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route for admin login of bo
	router.HandleFunc("/api/bo/users/login", func(w http.ResponseWriter, r *http.Request) {
		api.LoginAdmin(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route for creating a new role
	router.HandleFunc("/api/users/role", func(w http.ResponseWriter, r *http.Request) {
		api.CreateRole(client, dbName, os.Getenv("ROLES_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route to get technicians
	router.HandleFunc("/api/users/technicians", func(w http.ResponseWriter, r *http.Request) {
		api.GetTechnicians(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get clients
	router.HandleFunc("/api/users/clients", func(w http.ResponseWriter, r *http.Request) {
		api.GetClients(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get clients
	router.HandleFunc("/api/users/recovery-confirmation", func(w http.ResponseWriter, r *http.Request) {
		api.ConfirmAuthCode(client, dbName, userCollection, w, r)
	}).Methods("POST")
}
