package main

import (
	"chatter/trace"
	"context"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"github.com/stretchr/objx"
)

const (
	socketBufferSize = 1024 // socketBufferSize is a buffer size for reading and writing from the socket connection.
)

// upgrader is an instance of gorilla/websocket.Upgrader object use for upgrading
// the existing HTTP connection into a socket.
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

// newRoom is an helper function for creating a new room.
func newRoom(ctx context.Context, t trace.Tracer, a Avatar) *room {
	return &room{
		forward: make(chan *message),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  t,
		avatar: a,
		ctx:     ctx,
	}
}

type room struct {
	// forward is a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan *message
	// join is a channel for clients wishing to join the room.
	join chan *client
	// leave is a channel for clients wishing to leave the room.
	leave chan *client
	// clients holds all current clients in this room.
	clients map[*client]bool
	// tracer will receive trace information of activity.
	tracer trace.Tracer
	// avatar will be the getter of the avatar URL of the client.
	avatar Avatar
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
			r.tracer.Trace("New client joined")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client left")
		case msg := <-r.forward:
			r.tracer.Trace("Message received: ", msg.Message)
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		case <-r.ctx.Done():
			r.tracer.Trace("Room shutting down")
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
	authCookie, err := req.Cookie("auth")
	if err != nil {
		log.Fatal("chat: ", err.Error())
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan *message),
		room:   r,
		userData: objx.MustFromBase64(authCookie.Value),
	}
	r.join <- client
	defer func() {
		r.leave <- client
	}()
	go client.write()
	client.read()
}
