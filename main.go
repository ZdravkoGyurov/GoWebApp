package main

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	db "./db"
	handlers "./handlers"
)

func handleRoutes(collection *mongo.Collection) {
	dbc := handlers.DBCollection{collection}
	http.HandleFunc("/note", dbc.HandleNote)
	http.HandleFunc("/notes", dbc.HandleNotes)
}

func main() {
	client := db.Connect()
	defer db.Disconnect(client)

	collection := client.Database("test").Collection("notes")

	handleRoutes(collection)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
