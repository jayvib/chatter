package main

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"crypto/md5"
	"io"
	gomniauthcommon "github.com/stretchr/gomniauth/common"
	"log"
)

// ChatUser represent a single user of chatter this provide two methods
// that will be use whenever a function needs a unique id and avatar url.
type ChatUser interface {
	UniqueID() string
	AvatarURL() string
}

// chatUser is an implementation of ChatUser interface this can be
// a default ChatUser interface implementation if the client doesn't
// have its own implementation.
type chatUser struct {
	gomniauthcommon.User
	uniqueID string
}

// UniqueID provides a unique id of the chatter user.
func (u chatUser) UniqueID() string {
	return u.uniqueID
}

// newExternalAuth initialize the gomniauth providers that
// will be use during the OAUTH authentication of the user.
func newExternalAuth(conf *Config) {
	gomniauth.SetSecurityKey(conf.Auth.Facebook.Key)
	gomniauth.WithProviders(
		facebook.New(
			conf.Auth.Facebook.Key,
			conf.Auth.Facebook.Secret,
			conf.Auth.Facebook.URL,
		),
		google.New(
			conf.Auth.Google.Key,
			conf.Auth.Google.Secret,
			conf.Auth.Google.URL,
		),
	)
}

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
	cookie, err := r.Cookie("auth")
	if err != nil {
		if err == http.ErrNoCookie || cookie.Value == ""{
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
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusBadRequest)
			return
		}
		loginUrl, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to get auth url %s: %s", provider, err),
				http.StatusBadRequest)
			return
		}
		w.Header().Set("Location", loginUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to get provider %s: %s", provider, err),
				http.StatusBadRequest)
			return
		}
		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			http.Error(
				w,
				fmt.Sprintf("Error when trying to complete auth for" +
					"%s: %s", provider, err),
				http.StatusInternalServerError,
			)
			return
		}
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalln("Error when trying to get user from", provider, "-", err)
		}
		chatUser := &chatUser{ User: user }
		// create a user id using md5 hashing
		m := md5.New()
		io.WriteString(m, strings.ToLower(user.Email()))
		chatUser.uniqueID = fmt.Sprintf("%x", m.Sum(nil))
		avatarURL, err := avatars.GetAvatarURL(chatUser)
		if err != nil {
			log.Fatalln("Error when trying to GetAvatarURL", "-", err)
		}
		authCookieValue := objx.New(map[string]interface{}{
			"userid": chatUser.uniqueID,
			"name": user.Name(),
			"avatar_url": avatarURL,
			"email": user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}