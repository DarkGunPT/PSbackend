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

1. **PUT /api/v1/mb/users/register-completion:** Completes the registration process for Mobile App users.
    * Request body:
      ```json
      {
        "email": string,
        "name": string,
        "nif": string,
        "phone": string,
        "service_types": [string],
        "locality": string,
        "workStart": string,
        "workEnd": string
      }
    ```
1. **GET /api/v1/mb/users:** Retrieves all users for Mobile App.
1. **GET /api/v1/bo/users:** Retrieves all users for Back Office.

1. **GET /api/v1/mb/users/technicians:** Retrieves all technicians for Mobile App.
1. **GET /api/v1/bo/users/technicians:** Retrieves all technicians for Back Office.

1. **GET /api/v1/bo/users/nif:** Retrieves a user by NIF for Back Office.
    * Request body:
      ```json
      {
        "nif": int
      }
      ```
1. **PUT /api/v1/mb/users/{nif}:** Updates a user for Mobile App.
    * Request body:
      ```json
      {        
        "name": string,
        "password": string,
        "phone": int,
        "role": [string],
        "service_types": [string],
        "email": string,        
        "locality": string,
        "workStart": string,
        "workEnd": string
      }
      ```
1. **PUT /api/v1/bo/users/active:** Changes the `isActive` status of a user for Back Office.
    * Request body:
      ```json
      {
        "email": string
      }
      ```
1. **PUT /api/v1/bo/users/block:** Changes the `BlockServices` status of a user for Back Office.
    * Request body:
      ```json
      {
        "email": string
      }
      ```
1. **DELETE /api/v1/mb/users:** Deletes a user by NIF for Mobile App.
1. **DELETE /api/v1/bo/users:** Deletes a user by NIF for Back Office.
    * Request body:
      ```json
      {
        "nif": int
      }
      ```
1. **POST /api/v1/mb/users/login:** User login for Mobile App.
1. **POST /api/v1/bo/users/login:** Admin login for Back Office.
    * Request body:
      ```json
      {
        "email": string,
        "password": string
      }
      ```
1. **POST /api/v1/bo/users/role:** Creates a new role for Back Office.
    * Request body:
      ```json
      {
        "name": string
      }
      ```
1. **GET /api/v1/mb/users/clients:** Retrieves all clients for Mobile App.
1. **GET /api/v1/bo/users/clients:** Retrieves all clients for Back Office.

1. **POST /api/v1/mb/users/register:** Sends an email with a verification code for Mobile App registration.
1. **POST /api/v1/bo/users/register:** Sends an email with a verification code for Back Office registration.
    * Request body:
      ```json
      {
        "email": string
      }
      ```
1. **POST /api/v1/mb/users/register-confirmation:** Confirms the code and sets a new password for Mobile App registration.
1. **POST /api/v1/bo/users/register-confirmation:** Confirms the code and sets a new password for Back Office registration.
    * Request body:
      ```json
      {
        "email": string,
        "code": int,
        "password": string
      }
      ```
1. **POST /api/v1/mb/users/recovery:** Sends a recovery email with a verification code for Mobile App.
1. **POST /api/v1/bo/users/recovery:** Sends a recovery email with a verification code for Back Office.
    * Request body:
      ```json
      {
        "email": string,
      }
      ```
1. **POST /api/v1/mb/users/recovery-confirmation:** Confirms the recovery code and sets a new password for Mobile App.
1. **POST /api/v1/bo/users/recovery-confirmation:** Confirms the recovery code and sets a new password for Back Office.
    * Request body:
      ```json
      {
        "email": string,
        "code": int,
        "password": string
      }
      ```
1. **GET /api/v1/mb/users/{nif}:** Retrieves a user by NIF for Mobile App.

1. **GET /api/v1/bo/users/clients/order:** Retrieves clients ordered by a filter for Back Office.
1. **GET /api/v1/bo/users/technicians/order:** Retrieves technicians ordered by a filter for Back Office.
    * Request body:
      ```json
      {
        "filter": "rating" or "filter": "services"
      }
      ```
1. **GET /api/v1/bo/fees:** Retrieves all fees for Back Office.
    
1. **POST /api/v1/bo/fees:** Creates a fee for Back Office.
    * Request body:
      ```json
      {
        "nif": int,
        "value": float,
        "day": string,
        "month": string,
        "year": string
      }
1. **GET /api/v1/mb/fees/{nif}:** Retrieves fees of a technician by NIF for Mobile App.
1. **PUT /api/v1/mb/fees/{id}:** Updates the status of a fee to PAID for Mobile App.

### Services

1. **GET /api/v1/bo/services:** Retrieves all services for Back Office.
1. **GET /api/v1/mb/services:** Retrieves all services for Mobile App.

1. **GET /api/v1/bo/services/id:** Retrieves a service by ID for Back Office.
1. **GET /api/v1/mb/services/id:** Retrieves a service by ID for Mobile App.

1. **GET /api/v1/bo/services/service-type:** Retrieves all filtered services by type for Back Office.
1. **GET /api/v1/mb/services/service-type:** Retrieves all filtered services by type for Mobile App.
    * Request body:
      ```json
      {
        "name": string
      }
      ```
1. **PUT /api/v1/bo/services:** Updates a service for Back Office.
1. **PUT /api/v1/mb/services:** Updates a service for Mobile App.
    * Request body:
      ```json
      {
        "price": string or "name": string
      }
      ```
1. **POST /api/v1/bo/service-type:** Creates a new specific service type for Back Office.
    * Request body:
      ```json
      {
        "name": string
      }
      ```
1. **GET /api/v1/bo/service-type:** Retrieves all service types for Back Office.
1. **GET /api/v1/mb/service-type:** Retrieves all service types for Mobile App.

1. **PUT /api/v1/bo/service-type:** Updates a service type for Back Office.
    * Request body:
      ```json
      {
        "id": string,
        "name": string
      }
      ```
1. **DELETE /api/v1/bo/service-type:** Deletes a service type by ID for Back Office.
    * Request body:
      ```json
      {
        "id": string
      }
      ```
1. **GET /api/v1/mb/services/technicians:** Retrieves services by technician for Mobile App.
1. **GET /api/v1/bo/services/technicians:** Retrieves services by technician for Back Office.
    * Request body:
      ```json
      {
        "employee_id": string
      }
      ```
1. **POST /api/v1/mb/services/appointment:** Updates service with a new appointment for Mobile App.
    * Request body:
      ```json
      {
        "client_email": string,
        "provider_email": string,
        "service_name": string,
        "start": string,
        "end": string,
        "phone": string,
        "nif": string,
        "locality": string,
        "notes": string,
        "totalPrice": string
      }
      ```
1. **GET /api/v1/bo/services/appointments:** Retrieves appointments for Back Office.

1. **GET /api/v1/bo/services/appointments/upcoming:** Retrieves all upcoming appointments of a Technician.

1. **GET /api/v1/mb/services/appointments/upcoming/client/{nif}:** Retrieves all upcoming appointments of a Client.

1. **GET /api/v1/mb/services/appointments/upcoming/technician/{nif}:** Retrieves all upcoming appointments of a Technician.

1. **GET /api/v1/bo/services/appointments/history:** Retrieves history of appointments.

1. **GET /api/v1/mb/services/appointments/history/client/{nif}:** Retrieves history of appointments of a Client.
1. **GET /api/v1/mb/services/appointments/history/technician/{nif}:** Retrieves history of appointments of a Technician.

1. **GET /api/v1/bo/services/appointments/price:** Retrieves appointments in a price range.
    * Request body:
      ```json
      {
        "service_type": string,
        "max": float,
        "min": float
      }
      ```
1. **GET /api/v1/bo/services/price:** Retrieves services in a price range.
    * Request body:
      ```json
      {
        "service_type": string,
        "max": float,
        "min": float
      }
      ```
1. **DELETE /api/v1/mb/services/appointments/{id}:** Deletes an appointment by ID for Mobile App.