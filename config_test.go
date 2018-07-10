package main

import (
	"testing"
	"os"
)

func TestConfig(t *testing.T) {
	file, err := os.Open("config.json")
	if err != nil {
		t.Fatal(err)
	}
	conf, err := newConfig(file)
	if err != nil {
		t.Fatal("config must not be nil")
	}
	if conf.Auth.Facebook.Key != "191995201517761" {
		t.Errorf("expected facebook key doesn't matched with the actual key %s\n", conf.Auth.Facebook.Key)
	}
}
