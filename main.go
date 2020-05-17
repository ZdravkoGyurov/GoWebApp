package main

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ZdravkoGyurov/go-web-app/db"
	"github.com/ZdravkoGyurov/go-web-app/handlers"
)

func handleRoutes(collection *mongo.Collection) {
	http.HandleFunc("/note", handlers.HandleNote(collection))
	http.HandleFunc("/notes", handlers.HandleNotes(collection))
}

func main() {
	client := db.Connect()
	defer db.Disconnect(client)

	notesCollection := client.Database("go-web-app").Collection("notes")

	handleRoutes(notesCollection)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
