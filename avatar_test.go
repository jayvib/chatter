package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("AuthAvatar.GetAvatarURL should return" +
			"ErrNoAvatarURL when no value present")
	}
	testUrl := "http://url-to-gravatar/"
	client.userData = map[string]interface{}{
		"avatar_url": testUrl,
	}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
	}
	if url != testUrl {
		t.Error("URL got wasn't match from the URL data in the client")
	}
}

func TestGravatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{
		"email": "example.one@gmail.com",
	}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return an error")
	}
	expectedUrl := "//www.gravatar.com/avatar/2f64808765d95b66da158110bd756230"
	if url != expectedUrl {
		t.Errorf("GravatarAvatar.GetAvatarURL wrongly returned %s", url)
	}
}
