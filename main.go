package main

import (
	"net/http"
	"log"
	"context"
)



func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	r := newRoom(ctx)
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	go r.run()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
