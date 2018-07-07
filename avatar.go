package main

import "errors"

var ErrNoAvatarURL = errors.New("chatter: unable to get an avatar URL")

type Avatar interface {
	// GetAvatarURL will return the URL of the avatar
	// from any object where the avatar URL is get able
	// just by implementing the Avatar interface.
	//
	// It will return an ErrNoAvatarURL error if the
	// url wasn't found.
	GetAvatarURL(c *client) (url string, err error)
}