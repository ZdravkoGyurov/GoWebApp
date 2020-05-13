package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	models "./models"
)

// DBCollection encapsulates mongoDB collection
type DBCollection struct {
	collection *mongo.Collection
}

func parseTitle(url *url.URL) string {
	titles, ok := url.Query()["title"]

	if !ok || len(titles[0]) < 1 {
		log.Println("Title parameter is missing")
		return ""
	}

	return titles[0]
}

func (dbc DBCollection) readNote(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/note" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		title := parseTitle(r.URL)

		note := models.FindNoteByTitle(dbc.collection, title)
		fmt.Fprintln(w, note)
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error while parsing body: %v", err)
			return
		}
		title := r.Form.Get("title")
		body := r.Form.Get("body")

		models.InsertNote(dbc.collection, models.NewNote(title, body, time.Now()))
		fmt.Fprintln(w, "Insertion successful")
	case "PUT":
		r.ParseForm()
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "Error while parsing body: %v", err)
			return
		}

		newTitle := r.Form.Get("title")
		body := r.Form.Get("body")

		title := parseTitle(r.URL)
		models.UpdateNoteByTitle(dbc.collection, title, models.NewNote(newTitle, body, time.Now()))
		fmt.Fprintln(w, "Update successful")
	case "DELETE":
		title := parseTitle(r.URL)

		models.DeleteNoteByTitle(dbc.collection, title)
		fmt.Fprintln(w, "Delete successful")
	default:
		fmt.Fprintf(w, "Sorry, only GET, POST, PUT, DELETE methods are supported.")
	}
}

func (dbc DBCollection) readNotes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notes" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		notes := models.FindNotes(dbc.collection)
		fmt.Fprintln(w, notes)
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}

func handleRoutes(collection *mongo.Collection) {
	dbc := DBCollection{collection}
	http.HandleFunc("/note", dbc.readNote)
	http.HandleFunc("/notes", dbc.readNotes)
}

func main() {
	client := models.DbConnect()
	defer models.DbDisconnect(client)

	collection := client.Database("test").Collection("notes")

	handleRoutes(collection)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
