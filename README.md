# Microservices Architecture Project

  

![Microservices](https://img.shields.io/badge/Architecture-Microservices-blue)

![Go](https://img.shields.io/badge/Language-Go-00ADD8)

![Docker](https://img.shields.io/badge/Container-Docker-2496ED)

  

A modern, scalable microservices application built with Go, featuring multiple services that communicate with each other through a central API Gateway. This project demonstrates microservices best practices, containerization, and database integration.

  

## 📋 Table of Contents

  

- [Technology Stack](#technology-stack)

- [Architecture Overview](#architecture-overview)

- [Services](#services)

- [Project Structure](#project-structure)

- [Prerequisites](#prerequisites)

- [Installation & Setup](#installation--setup)

- [API Documentation](#api-documentation)

- [Testing](#testing)

- [Development](#development)

- [Troubleshooting](#troubleshooting)

- [Contributing](#contributing)

- [License](#license)

  

## 🔧 Technology Stack

  

### Backend

  

-  **Go** (Golang) - Primary programming language for all microservices

-  **Gin Framework** - HTTP web framework for the API Gateway and services

-  **GORM** - ORM library for database operations

  

### Databases

  

-  **PostgreSQL** - Relational database for User Service

-  **MongoDB** - NoSQL database for Product Service

  

### DevOps & Infrastructure

  

-  **Docker** - Containerization platform

-  **Docker Compose** - Multi-container orchestration

-  **Environment Configuration** - Using .env files

  

### Communication

  

-  **RESTful APIs** - For service-to-service communication

-  **JSON** - Data exchange format
## 📐 Architecture Overview

This project implements a microservices architecture with the following components:

```
            Client
              |
              v
           Gateway
              |
     +----------------+
     |                |
     v                v
User Service    Product Service
     |                |
     v                v
 PostgreSQL        MongoDB
```


  

### Components Overview

  

-  **Client:** Sends requests to the system.

-  **Gateway:** Acts as the entry point and routes requests to the appropriate microservices.

-  **User Service:** Handles user-related operations and stores data in PostgreSQL.

-  **Product Service:** Manages product-related operations and stores data in MongoDB.

  

## Services

  

### User Service

  

-  **Language**: Go

-  **Database**: PostgreSQL

-  **Description**: Manages user data and authentication

  

### Product Service

  

-  **Language**: Go

-  **Database**: MongoDB

-  **Description**: Manages product data and inventory

  

### API Gateway

  

-  **Language**: Go

-  **Framework**: Gin

-  **Description**: Routes requests to the appropriate service

  

## Project Structure

  

```

/project-root

│

├── /api-gateway

│ ├── main.go

│ └── ...

│

├── /user-service

│ ├── main.go

│ └── ...

│

├── /product-service

│ ├── main.go

│ └── ...

│

├── docker-compose.yml

├── .env.example

└── README.md

```

  

## Prerequisites

  

- Docker

- Docker Compose

- Go (for local development)

  

## Installation & Setup

  

1.  **Create an Environment File**

Copy the contents of `.env.example` into a new file named `.env` in the project root and adjust values as needed.

  

2.  **Start with Docker Compose**

```bash

docker-compose up --build

```

This command will build and start your services and databases.

  

## API Documentation

  

### User Service

  

-  **POST /users**

  

- Description: Create a new user

- Request Body: `{ "username": "string", "password": "string" }`

- Response: `{ "id": "string", "username": "string" }`

  

-  **GET /users/{id}**

- Description: Get user details by ID

- Response: `{ "id": "string", "username": "string" }`

  

### Product Service

  

-  **POST /products**

  

- Description: Create a new product

- Request Body: `{ "name": "string", "price": "number" }`

- Response: `{ "id": "string", "name": "string", "price": "number" }`

  

-  **GET /products/{id}**

- Description: Get product details by ID

- Response: `{ "id": "string", "name": "string", "price": "number" }`

  

## Testing

  

You can send requests to user and product endpoints through the Gateway (http://localhost:8000):

  

- Example: `POST http://localhost:8000/users`

- Example: `POST http://localhost:8000/products`

  

## Development

  

To run the services locally without Docker:

  

1.  **User Service**

  

```bash

cd user-service

go run main.go

```

  

2.  **Product Service**

  

```bash

cd product-service

go run main.go

```

  

3.  **API Gateway**

```bash

cd api-gateway

go run main.go

```

  

## Troubleshooting

  

-  **Database Connection Issues**: Ensure your `.env` file has the correct database connection strings.

-  **Service Not Starting**: Check the logs for error messages and ensure all dependencies are installed.

  

## Contributing

  

Contributions via pull requests or issues are welcome.

  

## License

  

This project is licensed under the MIT License.