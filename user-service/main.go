package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents a user entity in the database
type User struct {
    gorm.Model
    Name  string
    Email string
}

// registerServiceWithConsul registers this service with Consul for service discovery
func registerServiceWithConsul() *api.Client {
    config := api.DefaultConfig()
    config.Address = "consul:8500"

    client, err := api.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to connect to Consul: %v", err)
    }

    serviceID := "user-service"
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    "user-service",
        Address: "user-service",
        Port:    8001,
        Check: &api.AgentServiceCheck{
            HTTP:     "http://user-service:8001/health",
            Interval: "10s",
            Timeout:  "1s",
        },
    }

    err = client.Agent().ServiceRegister(registration)
    if err != nil {
        log.Fatalf("Error registering service with Consul: %v", err)
    }

    return client
}

func main() {
    // PostgreSQL connection information
    dbHost := os.Getenv("DB_HOST")
    dbUser := os.Getenv("DB_USER")
    dbPass := os.Getenv("DB_PASS")
    dbName := os.Getenv("DB_NAME")
    dbPort := os.Getenv("DB_PORT")

    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        dbHost, dbUser, dbPass, dbName, dbPort)

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    // Auto-create tables with GORM
    db.AutoMigrate(&User{})

    consulClient := registerServiceWithConsul()
    defer consulClient.Agent().ServiceDeregister("user-service")

    r := gin.Default()
    r.RedirectTrailingSlash = false

    // Health endpoint for service discovery checks
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "up"})
    })

    // CREATE - Add a new user
    r.POST("/users/", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        db.Create(&user)
        c.JSON(http.StatusOK, user)
    })

    // READ - Get all users
    r.GET("/users/", func(c *gin.Context) {
        var users []User
        db.Find(&users)
        c.JSON(http.StatusOK, users)
    })

    // READ - Get a single user by ID
    r.GET("/users/:id", func(c *gin.Context) {
        var user User
        if err := db.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        c.JSON(http.StatusOK, user)
    })

    // UPDATE - Update a user's information
    r.PUT("/users/:id", func(c *gin.Context) {
        var user User
        if err := db.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }

        var input User
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db.Model(&user).Updates(input)
        c.JSON(http.StatusOK, user)
    })

    // DELETE - Remove a user
    r.DELETE("/users/:id", func(c *gin.Context) {
        var user User
        if err := db.First(&user, c.Param("id")).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
            return
        }
        db.Delete(&user)
        c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
    })

    r.Run(":8001")
}