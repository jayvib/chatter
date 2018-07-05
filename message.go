package main

import "time"

// message represent the single message from the client
type message struct {
	Name    string
	Message string
	When    time.Time
	AvatarURL string
}
