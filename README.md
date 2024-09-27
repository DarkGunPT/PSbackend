# PSbackend

## Overview
This project is a REST API built with Go (Golang) and MongoDB. It provides endpoints for managing users and services. The API is designed to handle requests from a mobile application and a web application, allowing for CRUD operations on the database.

## Features
- **User Management:**
  - Create, read, update, and delete users.
  - User login with password verification.
- **Service Management:**
  - Create, read, update, and delete services.
  - Retrieve services by employee or employer ID.

## Technologies Used
- **Backend:** Go (Golang)
- **Database:** MongoDB

## Prerequisites
- A `.env` file with the following content: 

## API Endpoints

### Users
- POST /users: Create a new user
- GET /users: Fetch all users
- GET /users/{id}: Fetch a specific user by ID
- PUT /users/{id}: Update a specific user by ID
- DELETE /users/{id}: Delete a specific user by ID
- POST /login: Authenticate a user

### Services
- POST /services: Create a new service
- GET /services: Fetch all services
- GET /services/employee/{employee_id}: Fetch services by Employee ID
- GET /services/employer/{employer_id}: Fetch services by Employer ID
PUT /services/{id}: Update a specific service by ID
- DELETE /services/{id}: Delete a specific service by ID

## How to run 
1. git clone https://github.com/DarkGunPT/PSfrontend.git
2. cd ../PSbackend-directory
3. go mod tidy
4. go run main.go