package routes

import (
	"PSbackend/api"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoutes(client *mongo.Client, dbName, userCollection string, router *mux.Router) {
	// Define route to finish the registration for mobile
	router.HandleFunc("/api/v1/mb/users/register-completion", func(w http.ResponseWriter, r *http.Request) {
		api.RegisterCompletion(client, dbName, userCollection, w, r)
	}).Methods("PUT")

	// Define route to get users for mobile
	router.HandleFunc("/api/v1/mb/users", func(w http.ResponseWriter, r *http.Request) {
		api.GetUsers(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get users for backoffice
	router.HandleFunc("/api/v1/bo/users", func(w http.ResponseWriter, r *http.Request) {
		api.GetUsers(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get technicians for mobile
	router.HandleFunc("/api/v1/mb/users/technicians", func(w http.ResponseWriter, r *http.Request) {
		api.GetTechnicians(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get technicians for backoffice
	router.HandleFunc("/api/v1/bo/users/technicians", func(w http.ResponseWriter, r *http.Request) {
		api.GetTechnicians(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get a user by nif for backoffice
	router.HandleFunc("/api/v1/bo/users/nif", func(w http.ResponseWriter, r *http.Request) {
		api.GetUser(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to update a user for mobile
	router.HandleFunc("/api/v1/mb/users/{nif}", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateUser(client, dbName, userCollection, w, r)
	}).Methods("PUT")

	// Define route to change the isActive of a user
	router.HandleFunc("/api/v1/bo/users/active", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateActive(client, dbName, userCollection, w, r)
	}).Methods("PUT")

	// Define route to change the BlockServices of a user
	router.HandleFunc("/api/v1/bo/users/block", func(w http.ResponseWriter, r *http.Request) {
		api.UpdateBlock(client, dbName, userCollection, w, r)
	}).Methods("PUT")

	// Define route to delete a user by nif for mobile
	router.HandleFunc("/api/v1/mb/users", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteUser(client, dbName, userCollection, w, r)
	}).Methods("DELETE")

	// Define route to delete a user by nif for backoffice
	router.HandleFunc("/api/v1/bo/users", func(w http.ResponseWriter, r *http.Request) {
		api.DeleteUser(client, dbName, userCollection, w, r)
	}).Methods("DELETE")

	// Define route for user login of mobile
	router.HandleFunc("/api/v1/mb/users/login", func(w http.ResponseWriter, r *http.Request) {
		api.Login(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route for admin login of backoffice
	router.HandleFunc("/api/v1/bo/users/login", func(w http.ResponseWriter, r *http.Request) {
		api.LoginAdmin(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route for creating a new role for backoffice
	router.HandleFunc("/api/v1/bo/users/role", func(w http.ResponseWriter, r *http.Request) {
		api.CreateRole(client, dbName, os.Getenv("ROLES_COLLECTION"), w, r)
	}).Methods("POST")

	// Define route to get clients for mobile
	router.HandleFunc("/api/v1/mb/users/clients", func(w http.ResponseWriter, r *http.Request) {
		api.GetClients(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get clients for backoffice
	router.HandleFunc("/api/v1/bo/users/clients", func(w http.ResponseWriter, r *http.Request) {
		api.GetClients(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to send email with a code to verify for mobile
	router.HandleFunc("/api/v1/mb/users/register", func(w http.ResponseWriter, r *http.Request) {
		api.VerificateEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route to send email with a code to verify for backoffice
	router.HandleFunc("/api/v1/bo/users/register", func(w http.ResponseWriter, r *http.Request) {
		api.VerificateEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route to confirm the code set by the user and define the new password for mobile
	router.HandleFunc("/api/v1/mb/users/register-confirmation", func(w http.ResponseWriter, r *http.Request) {
		api.ConfirmAuthCode(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to confirm the code set by the user and define the new password for backoffice
	router.HandleFunc("/api/v1/bo/users/register-confirmation", func(w http.ResponseWriter, r *http.Request) {
		api.ConfirmAuthCode(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to send recovery email with a code to verify for mobile
	router.HandleFunc("/api/v1/mb/users/recovery", func(w http.ResponseWriter, r *http.Request) {
		api.RecoveryEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route to send recovery email with a code to verify for backoffice
	router.HandleFunc("/api/v1/bo/users/recovery", func(w http.ResponseWriter, r *http.Request) {
		api.RecoveryEmail(w, r, client, dbName, userCollection)
	}).Methods("POST")

	// Define route to confirm the code set by the user and define the new password for mobile
	router.HandleFunc("/api/v1/mb/users/recovery-confirmation", func(w http.ResponseWriter, r *http.Request) {
		api.ConfirmAuthCode(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to confirm the code set by the user and define the new password for backoffice
	router.HandleFunc("/api/v1/bo/users/recovery-confirmation", func(w http.ResponseWriter, r *http.Request) {
		api.ConfirmAuthCode(client, dbName, userCollection, w, r)
	}).Methods("POST")

	// Define route to get a user by nif  for mobile
	router.HandleFunc("/api/v1/mb/users/{nif}", func(w http.ResponseWriter, r *http.Request) {
		api.GetUser(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get clients ordened by a filter
	router.HandleFunc("/api/v1/bo/users/clients/order", func(w http.ResponseWriter, r *http.Request) {
		api.OrderClients(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get technicians ordened by a filter
	router.HandleFunc("/api/v1/bo/users/technicians/order", func(w http.ResponseWriter, r *http.Request) {
		api.OrderTechnicians(client, dbName, userCollection, w, r)
	}).Methods("GET")

	// Define route to get fees
	router.HandleFunc("/api/v1/bo/fees", func(w http.ResponseWriter, r *http.Request) {
		api.GetFees(client, dbName, os.Getenv("FEES_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route to get fees
	router.HandleFunc("/api/v1/bo/fees", func(w http.ResponseWriter, r *http.Request) {
		api.CreateFee(client, dbName, os.Getenv("FEES_COLLECTION"), userCollection, w, r)
	}).Methods("POST")

	// Define route to update the status of a fee to PAID
	router.HandleFunc("/api/v1/mb/fees/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.PayFee(client, dbName, os.Getenv("FEES_COLLECTION"), userCollection, w, r)
	}).Methods("PUT")

	// Define route to get fees of a technician
	router.HandleFunc("/api/v1/mb/fees/{nif}", func(w http.ResponseWriter, r *http.Request) {
		api.GetFeesByNif(client, dbName, os.Getenv("FEES_COLLECTION"), w, r)
	}).Methods("GET")

	// Define route to get fees of a technician
	router.HandleFunc("/api/v1/bo/count-services-performed", func(w http.ResponseWriter, r *http.Request) {
		api.GetServicesPerformed(client, dbName, userCollection, w, r)
	}).Methods("GET")
}
