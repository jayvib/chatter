package main

import "net/http"

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
