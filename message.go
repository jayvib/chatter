package main

import "time"

// message represent the single message from the client
type message struct {
	Name    string    `json:"Name"`
	Message string    `json:"Message"`
	When    time.Time `json:"When"`
}
