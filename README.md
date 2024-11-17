# FixFinder Backend

This repository contains the backend code for FixFinder, a platform designed to connect users with technicians for various domestic services. The backend is built using Go and MongoDB, providing a robust and scalable foundation for the application. The API is designed to handle requests from a mobile application and a web application, allowing for CRUD operations on the database.

## Features

* **User Management:**
    * User registration and authentication
    * User profile management
    * Role-based access control (e.g., technicians, clients)
    * User recovery and verification

* **Service Management:**
    * Service creation and management
    * Service type categorization
    * Service appointments

## API Endpoints

### Users

* **POST /api/users:** Creates a new user.
    * Request body:
        ```json
        {
          "email": string,
          "password": string,
        }
        ```
* **GET /api/users:** Retrieves all users.
* **GET /api/users/phone:** Retrieves a user by phone number.
    * Request body:
        ```json
        {
          "phone": int
        }
        ```
* **PUT /api/users:** Updates an existing user.
    * Request body:
        ```json
        {
          "nif": int,
          "name": string,
          "phone": int,
          "email": string,
          "password": string,
          "locality": string,
          "role": [string],
          "serviceTypes": [string]
        }
        ```
* **DELETE /api/users:** Deletes a user by NIF.
    * Request body:
        ```json
        {
          "nif": int
        }
        ```
* **POST /api/users/login:** Logs in a user.
    * Request body:
        ```json
        {
          "email": string,
          "password": string
        }
        ```
* **POST /api/bo/users/login:** Logs in an admin user (for back-office).
    * Request body:
        ```json
        {
          "email": string,
          "password": string
        }
        ```
* **POST /api/users/role:** Creates a new user role.
    * Request body:
        ```json
        {
          "name": string
        }
        ```
* **GET /api/users/technicians:** Retrieves all technicians.
* **GET /api/users/clients:** Retrieves all clients.
* **POST /api/mb/users/email:** Sends a verification code to the user's email for recovery (mobile users).
    * Request body:
        ```json
        {
          "email": string
        }
        ```
* **POST /api/users/email-confirmation:** Confirms an authentication code for user recovery.
    * Request body:
        ```json
        {
          "email": string,
          "code": int,
          "password": string
        }
        ```

### Services

* **POST /bo/services:** Creates a new service for Back Office.
* **POST /mb/services:** Creates a new service for Mobile App.
    * Request body:
        ```json
        {
          "serviceType": {
            "name": string
          },
          "description": string
        }
        ```
* **GET /bo/services:** Retrieves all services for Back Office.
* **GET /mb/services:** Retrieves all services for Mobile App.

* **GET /bo/services/id:** Retrieves a service by ID for Back Office.
* **GET /mb/services/id:** Retrieves a service by ID for Mobile App.
    * Request body:
        ```json
        {
          "id": string
        }
        ```
* **GET /bo/services/service-type:** Retrieves all filtered services by type for the Back Office.
* **GET /mb/services/service-type:** Retrieves all filtered services by type for the Mobile App.
    * Request body:
        ```json
        {
          "service_type": string
        }
      ```
* **PUT /bo/services:** Updates an existing service for Back Office.
* **PUT /mb/services:** Updates an existing service for Mobile App.
    * Request body:
        ```json
        {
          "id": string,
          "serviceType": {
            "name": string
          },
          "description": string
        }
        ```
* **POST /bo/service-type:** Creates a new service type for Back Office.
    * Request body:
        ```json
        {
          "name": string
        }
        ```
* **GET /bo/service-type:** Retrieves all service types for Back Office.
* **GET /mb/service-type:** Retrieves all service types for Mobile App.
* **PUT /bo/service-type:** Updates an existing service type for Back Office.
    * Request body:
        ```json
        {
          "id": string,
          "name": string
        }
        ```
* **DELETE /bo/service-type:** Deletes a service type by ID for Back Office.
    * Request body:
        ```json
        {
          "id": string
        }
        ```

## How to run

1. git clone https://github.com/DarkGunPT/PSfrontend.git
2. cd ../PSbackend-directory
3. go mod tidy
4. go run main.go