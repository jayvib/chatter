package main

import (
	"net/http"
	"strings"
	"log"
	"fmt"
)

// MustAuth is an helper function for authentication the handler before use by the
// client.
func MustAuth(h http.Handler) http.Handler {
	return &authHandler{
		next: h,
	}
}

// authHandler is an Handler for authenticating the user
// before using the chatter application.
type authHandler struct {
	next http.Handler
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie {
			// not authenticated, need to login
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}
		// display for now the other error to the HTML
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	a.next.ServeHTTP(w, r)
}

// loginHandler is an HandlerFunc use for logging in to the chat app.
// It uses third party login process.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	sep := strings.Split(r.URL.Path, "/")
	action := sep[2]
	provider := sep[3]
	switch action {
	case "login":
		log.Println("TODO handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}