package main

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ZdravkoGyurov/GoWebApp/db"
	"github.com/ZdravkoGyurov/GoWebApp/handlers"
)

func handleRoutes(collection *mongo.Collection) {
	http.HandleFunc("/note", handlers.HandleNote(collection))
	http.HandleFunc("/notes", handlers.HandleNotes(collection))
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
