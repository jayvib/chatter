package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client represent a single chat user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which messages are sent
	send chan *message
	// room is the room this client is chatting in
	room *room
	// userData holds information about the user
	userData map[string]interface{}
}

// read reads a message from the socket then forward the message to the room.
// If error occurs read function will terminate.
func (c *client) read() {
	defer c.socket.Close()
	for {
		msg := new(message)
		err := c.socket.ReadJSON(msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		msg.AvatarURL, _ = c.room.avatar.GetAvatarURL(c) // this need to be optimized.
		c.room.forward <- msg
	}
}

// write writes a message received from the send channel and write it to the
// socket. If any error occurs while writing to the socket this function will
// exit.
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
}
