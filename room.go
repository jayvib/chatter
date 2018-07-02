package main

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
	"log"
)

const (
	socketBufferSize = 1024 // socketBufferSize is a buffer size for reading and writing from the socket connection.
	messageBufferSize = 256 // messageBufferSize is a buffer size for storing the message from the clients.
)

// upgrader is an instance of gorilla/websocket.Upgrader object use for upgrading
// the existing HTTP connection into a socket.
var upgrader = &websocket.Upgrader{
	ReadBufferSize: socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// newRoom is an helper function for creating a new room.
func newRoom(ctx context.Context) *room {
	return &room {
		forward: make(chan []byte, messageBufferSize),
		join: make(chan *client),
		leave: make(chan *client),
		clients: make(map[*client]bool),
		ctx: ctx,
	}
}

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clinets.
	forward chan []byte
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// context is the mechanism for cleaning up the room
	ctx context.Context
}

// run serves as the gateway of the chat app. It handles the client that will join,
// leave and the message that have been forwarded to the room from the client
// to send to the different clients.
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		case <-r.ctx.Done():
			return
		}
	}
}

// ServeHTTP is an handler for a new request from the client.
// This will create a new instance of a client and append it to the
// room.
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("chat: ", err)
		return
	}
	client := &client{
		socket: socket,
		send: make(chan []byte),
		room: r,
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}



