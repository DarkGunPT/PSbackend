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
* **POST /api/mb/users/recovery:** Sends a verification code to the user's email for recovery (mobile users).
    * Request body:
        ```json
        {
          "email": string
        }
        ```
* **POST /api/bo/users/verification:** Sends a verification code to the user's email for verification (back-office users).
    * Request body:
        ```json
        {
          "email": string
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
* **POST /api/users/recovery-confirmation:** Confirms an authentication code for user recovery.
    * Request body:
        ```json
        {
          "email": string,
          "code": int
        }
        ```

### Services

* **POST /api/services:** Creates a new service.
    * Request body:
        ```json
        {
          "serviceType": {
            "name": string
          },
          "description": string
        }
        ```
* **GET /api/services:** Retrieves all services.
* **GET /api/services/id:** Retrieves a service by ID.
    * Request body:
        ```json
        {
          "id": string
        }
        ```
* **DELETE /api/services:** Deletes a service by ID.
    * Request body:
        ```json
        {
          "id": string
        }
        ```
* **PUT /api/services:** Updates an existing service.
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
* **POST /api/bo/services:** Creates a new service type.
    * Request body:
        ```json
        {
          "name": string
        }
        ```
* **GET /services/type:** Retrieves all service types.
* **PUT /services/type:** Updates an existing service type.
    * Request body:
        ```json
        {
          "id": string,
          "name": string
        }
        ```
* **DELETE /services/type:** Deletes a service type by ID.
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