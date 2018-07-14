package main

import (
	"net/http"
	"log"
	"context"
	"flag"
	"chatter/trace"
	"os"
)

var (
	// addr is the address of the chatter URL
	addr string
	// set the active Avatar implementation
	avatars Avatar = TryAvatars{
		UseFileSystemAvatar,
		UseGravatarAvatar,
		UseAuthAvatar,
	}
)
const configPath = "config.json"

func init() {
	flag.StringVar(&addr, "addr", ":8180", "The address of the chat application")
}

func main() {
	addr := ":" + os.Getenv("PORT")
	if addr == ":" {
		log.Fatalln("chatter: port can't be found in environment variable")
	}
	flag.Parse()
	confFile, err := os.Open(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	conf, err := newConfig(confFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	newExternalAuth(conf)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//tracer := trace.New(os.Stdout)
	r := newRoom(ctx, trace.Off(), UseAuthAvatar)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/upload", &templateHandler{filename: "upload.html"})
	http.HandleFunc("/uploader", uploaderHandler)
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request){
		cookie := &http.Cookie{
			Name: "auth",
			Value: "",
			Path: "/",
			MaxAge: -1, // delete the cookie immediately
		}
		http.SetCookie(w, cookie)
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.Handle("/avatars/",
		http.StripPrefix("/avatars/",
			http.FileServer(http.Dir("./avatars"))))
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
