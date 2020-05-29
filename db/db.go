package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// DbUsername used to connect to MongoDB
	DbUsername = "zdravko"

	//DbPassword used to connect to Mongodb
	DbPassword = "pass"

	// DbHost is the MongoDB host / deployment name
	DbHost = "my-release-mongodb"

	// DbPort is the MongoDB port
	DbPort = "27017"

	// DbDefault is the initial database
	DbDefault = "go-web-app"

	// MongoURI is the whole MongoDB URI string including the username, password, host, port and initial database
	MongoURI = "mongodb://%s:%s@%s:%s/%s"
)

// Connect to MongoDB and return a client
func Connect() *mongo.Client {
	// clientOptions := options.Client().ApplyURI(fmt.Sprintf(MongoURI, DbUsername, DbPassword, DbHost, DbPort, DbDefault))
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println(err)
		log.Fatal(errors.New("An error occured while trying to connect to MongoDB"))
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println(err)
		log.Fatal(errors.New("MongoDB is not responding"))
	}

	fmt.Println("Connected to MongoDB!")

	return client
}

// Disconnect from MongoDB
func Disconnect(client *mongo.Client) {
	err := client.Disconnect(context.TODO())
	if err != nil {
		log.Println(err)
		log.Fatal("An error occured while trying to disconnect from MongoDB")
	}
	fmt.Println("Disconnected from MongoDB.")
}
