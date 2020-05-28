package main

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ZdravkoGyurov/go-web-app/db"
	"github.com/ZdravkoGyurov/go-web-app/handlers"
	"github.com/gorilla/sessions"
)

func handleRoutes(collection *mongo.Collection, store *sessions.CookieStore) {
	http.HandleFunc("/note", handlers.HandleNote(collection, store))
	http.HandleFunc("/notes", handlers.HandleNotes(collection, store))

	http.Handle("/", http.FileServer(http.Dir("templates/")))

	http.HandleFunc("/auth/login", handlers.OauthGoogleLogin(store))
	http.HandleFunc("/auth/callback", handlers.OauthGoogleCallback(store))
	http.HandleFunc("/auth/logout", handlers.Logout(store))
}

// os.Setenv("GOOGLE_OAUTH2_CLIENT_ID", "874806098368-u04ecnakpci7olgqrihb6af1978jrr4h.apps.googleusercontent.com")
// os.Setenv("GOOGLE_OAUTH2_CLIENT_SECRET", "DsqYhScPag1Y_NuvZI-laGwY")

func main() {
	client := db.Connect()
	defer db.Disconnect(client)

	// TODO: notesCollection := client.Database("go-web-app").Collection("notes")
	notesCollection := client.Database("test").Collection("notes")

	store := sessions.NewCookieStore([]byte("SESSION_KEY"))

	handleRoutes(notesCollection, store)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
