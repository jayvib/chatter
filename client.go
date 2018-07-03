package main

import "github.com/gorilla/websocket"

// client represent a single chat user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

// read reads a message from the socket then forward the message to the room.
// If error occurs read function will terminate.
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

// write writes a message received from the send channel and write it to the
// socket. If any error occurs while writing to the socket this function will
// exit.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
