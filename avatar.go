package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

var ErrNoAvatarURL = errors.New("chatter: unable to get an avatar URL")

type Avatar interface {
	// GetAvatarURL will return the URL of the avatar
	// from any object where the avatar URL is get able
	// just by implementing the Avatar interface.
	//
	// It will return an ErrNoAvatarURL error if the
	// url wasn't found.
	GetAvatarURL(c ChatUser) (url string, err error)
}

var UseAuthAvatar AuthAvatar

// AuthAvatar is an implementation of Avatar interface. This will get the avatar url
// provided by the third party OAUTH provider that is stored in client user data.
type AuthAvatar struct {}

func (AuthAvatar) GetAvatarURL(c ChatUser) (string, error) {
	url := c.AvatarURL()
	if url == "" {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

var UseGravatarAvatar GravatarAvatar

// GravatarAvatar is an implementation of Avatar interface. This will get the avatar url
// from the Gravatar api using the email provided by the client.
type GravatarAvatar struct {}

func (GravatarAvatar) GetAvatarURL(c ChatUser) (string, error) {
	if userId := c.UniqueID(); userId != "" {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userId), nil
	}
	return "", ErrNoAvatarURL
}

var UseFileSystemAvatar FileSystemAvatar

type FileSystemAvatar struct {}

func (FileSystemAvatar) GetAvatarURL(c ChatUser) (string, error) {
	if userId := c.UniqueID(); userId != "" {
		files, err := ioutil.ReadDir("avatars")
		if err != nil {
			return "", ErrNoAvatarURL
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			match, err := path.Match(userId + "*", file.Name())
			if err != nil {
				return "", ErrNoAvatarURL
			}
			if match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(c ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(c); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}