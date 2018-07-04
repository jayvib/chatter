package main

import (
	"net/http"
	"log"
	"context"
	"flag"
	"chatter/trace"
	"os"
)

var addr string

const configPath = "config.json"

func init() {
	flag.StringVar(&addr, "addr", ":8080", "The address of the chat application")
}

func main() {
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
	r := newRoom(ctx, trace.Off())
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
