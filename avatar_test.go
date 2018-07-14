package main

import (
	"testing"
	"path/filepath"
	"io/ioutil"
	"os"
	gomnitest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomnitest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{ User: testUser }
	url, err := authAvatar.GetAvatarURL(testChatUser)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return" +
			"ErrNoAvatarURL when no value present")
	}
	testUrl := "http://url-to-gravatar/"
	testUser = &gomnitest.TestUser{}
	testUser.On("AvatarURL").Return(testUrl, nil)
	testChatUser.User = testUser
	url, err = authAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("URL got wasn't match from the URL data in the client")
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	testChatUser := &chatUser{
		uniqueID: "2f64808765d95b66da158110bd756230",
	}
	url, err := gravatarAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	expectedUrl := "//www.gravatar.com/avatar/2f64808765d95b66da158110bd756230"
	if url != expectedUrl {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer os.Remove(filename)
	var fileSystemAvatar FileSystemAvatar
	testChatUser := &chatUser{
		uniqueID: "abc",
	}
	url, err := fileSystemAvatar.GetAvatarURL(testChatUser)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return an error")
	}
	if url != "/avatars/abc.jpg" {
		t.Error("received url wasn't matched from the expected url")
	}
}