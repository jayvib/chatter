package main

import (
	"net/http"
	"log"
	"context"
	"flag"
)

var addr string

func init() {
	flag.StringVar(&addr, "addr", ":8080", "The address of the chat application")
}

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := newRoom(ctx)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	log.Println("Starting web server on", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
