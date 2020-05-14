package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	models "../db/models"

	"go.mongodb.org/mongo-driver/mongo"
)

// DBCollection encapsulates mongoDB collection
type DBCollection struct {
	Collection *mongo.Collection
}

func parseTitle(url *url.URL) string {
	titles, ok := url.Query()["title"]

	if !ok || len(titles[0]) < 1 {
		log.Println("Title parameter is missing")
		return ""
	}

	return titles[0]
}

func getNote(dbc DBCollection, w http.ResponseWriter, r *http.Request) {
	title := parseTitle(r.URL)

	if note, err := models.FindNoteByTitle(dbc.Collection, title); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, note)
	}
}

func postNote(dbc DBCollection, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error while parsing body: %v", err)
		return
	}

	title := r.Form.Get("title")
	body := r.Form.Get("body")
	newNote := models.NewNote(title, body, time.Now())

	if err := models.InsertNote(dbc.Collection, newNote); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "Insertion successful")
	}
}

func putNote(dbc DBCollection, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error while parsing body: %v", err)
		return
	}

	title := parseTitle(r.URL)
	newTitle := r.Form.Get("title")
	body := r.Form.Get("body")
	newNote := models.NewNote(newTitle, body, time.Now())

	if err := models.UpdateNoteByTitle(dbc.Collection, title, newNote); err != nil {
		fmt.Fprintln(w, err, title)
		fmt.Println(err, title)
	} else {
		fmt.Fprintln(w, "Update successful")
	}
}

func deleteNote(dbc DBCollection, w http.ResponseWriter, r *http.Request) {
	title := parseTitle(r.URL)

	if err := models.DeleteNoteByTitle(dbc.Collection, title); err != nil {
		fmt.Fprintln(w, err, title)
		fmt.Println(err, title)
	} else {
		fmt.Fprintln(w, "Delete successful")
	}
}

func getNotes(dbc DBCollection, w http.ResponseWriter) {
	if notes, err := models.FindNotes(dbc.Collection); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, notes)
	}
}

// HandleNote handles GET, POST, PUT and DELETE requests to /note
func (dbc DBCollection) HandleNote(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/note" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		getNote(dbc, w, r)
	case "POST":
		postNote(dbc, w, r)
	case "PUT":
		putNote(dbc, w, r)
	case "DELETE":
		deleteNote(dbc, w, r)
	default:
		fmt.Fprintf(w, "Sorry, only GET, POST, PUT, DELETE methods are supported.")
	}
}

// HandleNotes handles GET request to /notes
func (dbc DBCollection) HandleNotes(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/notes" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		getNotes(dbc, w)
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}
