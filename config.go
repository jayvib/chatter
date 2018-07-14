package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

type Config struct {
	Auth struct {
		Facebook struct {
			Key    string `json:"key"`
			Secret string `json:"secret"`
			URL    string `json:"url"`
		} `json:"facebook"`
		Google struct {
			Key    string `json:"key"`
			Secret string `json:"secret"`
			URL    string `json:"url"`
		} `json:"google"`
		Github struct {
			Key    string `json:"key"`
			Secret string `json:"secret"`
			URL    string `json:"url"`
		} `json:"github"`
	} `json:"auth"`
}

func newConfig(confFile io.Reader) (*Config, error) {
	conf := new(Config)
	b, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
