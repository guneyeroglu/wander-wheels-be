# Wander Wheels Back-End Documentation

## Overview

Welcome to the Back-End repository of Wander Wheels, a rental cars application. This Back-End is built using Golang and several other technologies to ensure a robust and efficient service.

## Technologies Used

- **Golang**: The core programming language used for this project, known for its performance and efficiency.
- **Fiber**: An Express-inspired web framework for Golang that is designed for ease of use and high performance.
- **Viper**: A configuration management package that provides a comprehensive solution for Go applications.
- **pq**: A pure Go Postgres driver for the database interaction.
- **JWT**: JSON Web Tokens for secure user authentication.

## Getting Started

### Prerequisites

- Go (version 1.22.1)
- PostgreSQL (version 16+)

### Installation

1. **Clone the repository:**
   ```sh
   git clone https://github.com/guneyeroglu/wander-wheels-be.git
   cd wander-wheels-be
   ```

2. **Install dependencies:**
   ```sh
   go mod tidy
   ```

3. **Environment variables:**
   ##### Create a .env file in the root directory and add the necessary configuration variables. You can use the following template:
   ```sh
   DB_USERNAME = "yourUsername"
   DB_PASSWORD = "yourPassword"
   DB_CONNECTION_URL = "localhost"
   DB_PORT = "5432"
   DB_NAME = "databaseName"
   JWT_TOKEN_TYPE = "Bearer "
   JWT_SECRET_CODE = "yourSecretCode"
   ```

### Running the Application
  ```sh
  go run main.go
  ```
If you're using air you can simply run
  ```sh
  air
  ```

And now, you're ready to go.

## API Endpoints

### GET
-  /user-info
-  /fuels
-  /transmissions
-  /colors
-  /brands
-  /models
-  /cities
-  /cities/:id
-  /cars/:id
-  /price-range
-  /year-range
-  /seats

### POST
-  /login
-  /sign-up
-	 /cars
-	 /rent-car

#

This README was prepared by ChatGPT.
	
