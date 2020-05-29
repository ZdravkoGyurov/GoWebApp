package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH2_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleLogin(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectIfAlreadyLoggedIn(w, r, store)

		oauthState := generateStateOauthCookie(w)

		url := googleOauthConfig.AuthCodeURL(oauthState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}

}

func redirectIfAlreadyLoggedIn(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore) {
	sessionIDcookie, err := r.Cookie("session-id")
	if err == nil {
		http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
	}

	if sessionIDcookie != nil {
		session, err := store.Get(r, sessionIDcookie.Value)
		if err == nil || len(session.Values) != 0 {
			http.Redirect(w, r, "/notes", http.StatusTemporaryRedirect)
		}
	}
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)

	return state
}

func OauthGoogleCallback(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState, _ := r.Cookie("oauthstate")
		if r.FormValue("state") != oauthState.Value {
			log.Println("invalid oauth google state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		data, err := getUserDataFromGoogle(r.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		userInfo := UserInfo{}
		json.Unmarshal([]byte(data), &userInfo)

		sessionID, err := saveUserInfoInSession(w, r, store, userInfo)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		http.SetCookie(w, &http.Cookie{Name: "session-id", Value: sessionID, Path: "/"})

		http.Redirect(w, r, "/notes", http.StatusPermanentRedirect)
	}
}

func generateRandomSessionID() string {
	buff := make([]byte, 16)
	rand.Read(buff)
	return base64.RawURLEncoding.EncodeToString(buff)
}

func saveUserInfoInSession(w http.ResponseWriter, r *http.Request, store *sessions.CookieStore, userInfo UserInfo) (string, error) {
	sessionID := generateRandomSessionID()
	session, err := store.Get(r, sessionID)
	if err != nil {
		fmt.Fprintf(w, "Could not create session")
	}

	session.Options.MaxAge = 5 * 60
	session.Values["email"] = userInfo.Email
	session.Values["name"] = userInfo.Name

	err = session.Save(r, w)
	if err != nil {
		return "", fmt.Errorf("Could not save session in store: %s", err.Error())
	}

	return sessionID, nil
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	response, err := http.Get(oauthGoogleURLAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func Logout(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionIDcookie, err := r.Cookie("session-id")
		if err != nil {
			fmt.Printf("session ID missing: %s", err.Error())
			return
		}

		session, err := store.Get(r, sessionIDcookie.Value)
		if err != nil || len(session.Values) == 0 {
			fmt.Printf("could not get session: %s", err.Error())
			deleteCookie(w, sessionIDcookie.Name)
			return
		}

		deleteCookie(w, sessionIDcookie.Name)
		// TODO: session is not getting deleted
		session.Options.MaxAge = 1
		err = session.Save(r, w)
		if err != nil {
			fmt.Printf("could not delete session: %s", err.Error())
			return
		}

		http.Redirect(w, r, "/", http.StatusPermanentRedirect)
	}
}

func deleteCookie(w http.ResponseWriter, name string) {
	c := &http.Cookie{
		Name:   name,
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, c)
}
