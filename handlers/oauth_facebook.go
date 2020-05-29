package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var facebookOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/facebook/callback",
	ClientID:     os.Getenv("FACEBOOK_OAUTH2_CLIENT_ID"),
	ClientSecret: os.Getenv("FACEBOOK_OAUTH2_CLIENT_SECRET"),
	Scopes:       []string{"email"},
	Endpoint:     facebook.Endpoint,
}

const oauthFacebookURLAPI = "https://graph.facebook.com/me?access_token="

func OauthFacebookLogin(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		redirectIfAlreadyLoggedIn(w, r, store)

		oauthState := generateStateOauthCookie(w)
		url := facebookOauthConfig.AuthCodeURL(oauthState)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func OauthFacebookCallback(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oauthState, _ := r.Cookie("oauthstate")
		if r.FormValue("state") != oauthState.Value {
			log.Println("invalid oauth facebook state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		data, err := getUserDataFromFacebook(r.FormValue("code"))
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

func getUserDataFromFacebook(code string) ([]byte, error) {
	token, err := facebookOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	fields := "&fields=name,email"
	response, err := http.Get(oauthFacebookURLAPI + token.AccessToken + fields)
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
