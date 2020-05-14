package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoURI MongoDB connection URI
const MongoURI = "mongodb://localhost:27017"

// Connect to MongoDB and return a client
func Connect() *mongo.Client {
	clientOptions := options.Client().ApplyURI(MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(errors.New("An error occured while trying to connect to MongoDB"))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(errors.New("MongoDB is not responding"))
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

// Disconnect from MongoDB
func Disconnect(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal("An error occured while trying to disconnect from MongoDB")
	}
	fmt.Println("Disconnected from MongoDB.")
}
