package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/ZdravkoGyurov/go-web-app/db/models"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

func parseTitle(url *url.URL) string {
	titles, ok := url.Query()["title"]

	if !ok || len(titles[0]) < 1 {
		log.Println("Title parameter is missing")
		return ""
	}

	return titles[0]
}

func getNote(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	title := parseTitle(r.URL)

	if note, err := models.FindNoteByTitle(collection, title); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, note)
	}
}

func postNote(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error while parsing body: %v", err)
		return
	}

	title := r.Form.Get("title")
	body := r.Form.Get("body")
	newNote := models.NewNote(title, body, time.Now())

	if err := models.InsertNote(collection, newNote); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, "Insertion successful")
	}
}

func putNote(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error while parsing body: %v", err)
		return
	}

	title := parseTitle(r.URL)
	newTitle := r.Form.Get("title")
	body := r.Form.Get("body")
	newNote := models.NewNote(newTitle, body, time.Now())

	if err := models.UpdateNoteByTitle(collection, title, newNote); err != nil {
		fmt.Fprintln(w, err, title)
		fmt.Println(err, title)
	} else {
		fmt.Fprintln(w, "Update successful")
	}
}

func deleteNote(collection *mongo.Collection, w http.ResponseWriter, r *http.Request) {
	title := parseTitle(r.URL)

	if err := models.DeleteNoteByTitle(collection, title); err != nil {
		fmt.Fprintln(w, err, title)
		fmt.Println(err, title)
	} else {
		fmt.Fprintln(w, "Delete successful")
	}
}

func getNotes(collection *mongo.Collection, w http.ResponseWriter) {
	if notes, err := models.FindNotes(collection); err != nil {
		fmt.Fprintln(w, err)
		fmt.Println(err)
	} else {
		fmt.Fprintln(w, notes)
	}
}

// HandleNote handles GET, POST, PUT and DELETE requests to /note
func HandleNote(collection *mongo.Collection, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := validateUser(w, r, store)
		if err != nil {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		if r.URL.Path != "/note" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			email := session.Values["email"]
			name := session.Values["name"]
			fmt.Fprintf(w, "Hello, %s(%s)\n", name, email)

			getNote(collection, w, r)
		case "POST":
			postNote(collection, w, r)
		case "PUT":
			putNote(collection, w, r)
		case "DELETE":
			deleteNote(collection, w, r)
		default:
			fmt.Fprintf(w, "Sorry, only GET, POST, PUT, DELETE methods are supported.")
		}
	}
}

// HandleNotes handles GET request to /notes
func HandleNotes(collection *mongo.Collection, store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := validateUser(w, r, store)
		if err != nil {
			http.Error(w, "401 unauthorized", http.StatusUnauthorized)
			return
		}

		email := session.Values["email"]
		name := session.Values["name"]
		fmt.Fprintf(w, "Hello, %s(%s)\n", name, email)

		if r.URL.Path != "/notes" {
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}

		switch r.Method {
		case "GET":
			getNotes(collection, w)
		default:
			fmt.Fprintf(w, "Sorry, only GET method is supported.")
		}
	}
}

func validateUser(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) (*sessions.Session, error) {
	sessionIDcookie, err := r.Cookie("session-id")
	if err != nil {
		return nil, fmt.Errorf("session ID missing: %s", err.Error())
	}

	session, err := store.Get(r, sessionIDcookie.Value)
	if err != nil || len(session.Values) == 0 {
		return nil, fmt.Errorf("could not get session: %s", err.Error())
	}

	return session, nil
}
