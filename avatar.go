package main

import (
	"errors"
	"crypto/md5"
	"io"
	"fmt"
	"strings"
)

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

var UseAuthAvatar AuthAvatar

// AuthAvatar is an implementation of Avatar interface. This will get the avatar url
// provided by the third party OAUTH provider that is stored in client user data.
type AuthAvatar struct {}

func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

var UseGravatarAvatar GravatarAvatar

// GravatarAvatar is an implementation of Avatar interface. This will get the avatar url
// from the Gravatar api using the email provided by the client.
type GravatarAvatar struct {}

func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
