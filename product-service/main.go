package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Product represents a product item in the database
type Product struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name  string             `bson:"name,omitempty" json:"name,omitempty"`
    Price float64            `bson:"price,omitempty" json:"price,omitempty"`
}

var productCollection *mongo.Collection

// registerServiceWithConsul registers this service with Consul for service discovery
func registerServiceWithConsul() *api.Client {
    config := api.DefaultConfig()
    config.Address = "consul:8500"

    client, err := api.NewClient(config)
    if err != nil {
        log.Fatalf("Failed to connect to Consul: %v", err)
    }

    serviceID := "product-service"
    registration := &api.AgentServiceRegistration{
        ID:      serviceID,
        Name:    "product-service",
        Address: "product-service",
        Port:    8002,
        Check: &api.AgentServiceCheck{
            HTTP:     "http://product-service:8002/health",
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
    mongoHost := os.Getenv("MONGO_HOST")
    mongoPort := os.Getenv("MONGO_PORT")
    mongoDB := os.Getenv("MONGO_DB")

    // Construct MongoDB connection string
    mongoURI := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort)

    client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal("Failed to create MongoDB client:", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err = client.Connect(ctx); err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Initialize the products collection
    productCollection = client.Database(mongoDB).Collection("products")

    consulClient := registerServiceWithConsul()

    // Deregister service when application shuts down
    defer func() {
        consulClient.Agent().ServiceDeregister("product-service")
    }()

    r := gin.Default()
    r.RedirectTrailingSlash = false

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "product-service up"})
    })

    // CRUD endpoints
    r.POST("/products/", createProduct)
    r.GET("/products/", getAllProducts)
    r.GET("/products/:id", getProductByID)
    r.PUT("/products/:id", updateProduct)
    r.DELETE("/products/:id", deleteProduct)

    log.Fatal(r.Run(":8002"))
}

// createProduct handles the creation of a new product
func createProduct(c *gin.Context) {
    var product Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    res, err := productCollection.InsertOne(ctx, product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// getAllProducts retrieves all products from the database
func getAllProducts(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cursor, err := productCollection.Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    var products []Product
    if err := cursor.All(ctx, &products); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, products)
}

// getProductByID retrieves a specific product by its ID
func getProductByID(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var product Product
    err = productCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&product)
    if err == mongo.ErrNoDocuments {
        c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
        return
    } else if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, product)
}

// updateProduct updates an existing product
func updateProduct(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var product Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    update := bson.M{"$set": product}

    _, err = productCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// deleteProduct removes a product from the database
func deleteProduct(c *gin.Context) {
    id := c.Param("id")
    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    _, err = productCollection.DeleteOne(ctx, bson.M{"_id": objID})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}