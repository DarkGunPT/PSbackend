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

* **GET /api/v1/mb/users:** Retrieves all users for Mobile App.
* **GET /api/v1/bo/users:** Retrieves all users for Back Office.

* **GET /api/v1/mb/users/nif:** Retrieves a user by NIF for Mobile App.
* **GET /api/v1/bo/users/nif:** Retrieves a user by NIF for Back Office.
    * Request body:
        ```json
        {
          "nif": int
        }
        ```
* **PUT /api/v1/mb/users:** Updates a user for Mobile App.
* **PUT /api/v1/bo/users:** Updates a user for Back Office.
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
* **DELETE /api/v1/mb/users:** Deletes a user by NIF for Mobile App.
* **DELETE /api/v1/bo/users:** Deletes a user by NIF for Back Office.
    * Request body:
        ```json
        {
          "nif": int
        }
        ```
* **POST /api/v1/mb/users/login:** User login for Mobile App.
* **POST /api/v1/bo/users/login:** Admin login for Back Office.
    * Request body:
        ```json
        {
          "email": string,
          "password": string
        }
        ```
* **POST /api/v1/bo/users/role:** Creates a new role for Back Office.
    * Request body:
        ```json
        {
          "name": string
        }
        ```
* **GET /api/v1/mb/users/technicians:** Retrieves all technicians for Mobile App.
* **GET /api/v1/bo/users/technicians:** Retrieves all technicians for Back Office.
* **GET /api/v1/mb/users/clients:** Retrieves all clients for Mobile App.
* **GET /api/v1/bo/users/clients:** Retrieves all clients for Back Office.

* **POST /api/v1/mb/users/email:** Sends an email with a code to verify for Mobile App.
* **POST /api/v1/bo/users/email:** Sends an email with a code to verify for Back Office.
    * Request body:
        ```json
        {
          "email": string
        }
        ```
* **POST /api/v1/mb/users/email-confirmation:** Confirms the code and sets a new password for Mobile App.
* **POST /api/v1/bo/users/email-confirmation:** Confirms the code and sets a new password for Back Office.recovery.
    * Request body:
        ```json
        {
          "email": string,
          "code": int,
          "password": string
        }
        ```
* **POST /api/v1/mb/users/recovery:** Sends a recovery email with a code to verify for Mobile App.
* **POST /api/v1/bo/users/recovery:** Sends a recovery email with a code to verify for Back Office.
    * Request body:
        ```json
        {
          "email": string,
        }
        ```
* **POST /api/v1/mb/users/recovery-confirmation:** Confirms the recovery code and sets a new password for Mobile App.
* **POST /api/v1/bo/users/recovery-confirmation:** Confirms the recovery code and sets a new password for Back Office.
    * Request body:
        ```json
        {
          "email": string,
          "code": int,
          "password": string
        }
        ```
* **PUT /api/v1/mb/users/registration-confirmation:** Finishes the registration for Mobile App.
    * Request body:
      ```json
      {
        "name": string,
        "nif": int,
        "phone": int,
        "role": [string],
        "service_types": [string],
        "locality": string,
        "is_active": bool
      }
      ```

### Services

* **POST /api/v1/bo/services:** Creates a new service for Back Office.
* **POST /api/v1/mb/services:** Creates a new service for Mobile App.
    * Request body:
        ```json
        {
          "serviceType": {
            "name": string
          },
          "description": string
        }
        ```
* **GET /api/v1/bo/services:** Retrieves all services for Back Office.
* **GET /api/v1/mb/services:** Retrieves all services for Mobile App.

* **GET /api/v1/bo/services/id:** Retrieves a service by ID for Back Office.
* **GET /api/v1/mb/services/id:** Retrieves a service by ID for Mobile App.
    * Request body:
        ```json
        {
          "id": string
        }
        ```
* **GET /api/v1/bo/services/service-type:** Retrieves all filtered services by type for the Back Office.
* **GET /api/v1/mb/services/service-type:** Retrieves all filtered services by type for the Mobile App.
    * Request body:
        ```json
        {
          "service_type": string
        }
      ```
* **PUT /api/v1/bo/services:** Updates an existing service for Back Office.
* **PUT /api/v1/mb/services:** Updates an existing service for Mobile App.
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
* **POST /api/v1/bo/service-type:** Creates a new service type for Back Office.
    * Request body:
        ```json
        {
          "name": string
        }
        ```
* **GET /api/v1/bo/service-type:** Retrieves all service types for Back Office.
* **GET /api/v1/mb/service-type:** Retrieves all service types for Mobile App.
* **PUT /api/v1/bo/service-type:** Updates an existing service type for Back Office.
    * Request body:
        ```json
        {
          "id": string,
          "name": string
        }
        ```
* **DELETE /api/v1/bo/service-type:** Deletes a service type by ID for Back Office.
    * Request body:
        ```json
        {
          "id": string
        }
        ```
* **GET /api/v1/mb/services/technicians:** Retrieves all services by technician for the Mobile App.
* **GET /api/v1/bo/services/technicians:** Retrieves all services by technician for the Back Office.
    * Request body:
        ```json
        {
          "employee_id": string
        }
        ```

## How to run

1. git clone https://github.com/DarkGunPT/PSfrontend.git
2. cd ../PSbackend-directory
3. go mod tidy
4. go run main.go