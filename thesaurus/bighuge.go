package thesaurus

import (
	"net/http"
	"fmt"
	"encoding/json"
)

var url_format = "http://words.bighugelabs.com/api/2/%s/%s/json"

func New(apiKey string) Thesaurus {
	return &bigHuge{
		APIKey: apiKey,
	}
}

type bigHuge struct {
	APIKey string
}

type synonyms struct {
	Noun *words `json:"noun"`
	Verb *words `json:"verb"`
}

type words struct {
	Syn []string `json:"syn"`
}

func (b *bigHuge) Synonyms(term string) ([]string, error) {
	var syn []string
	resp, err := http.Get(fmt.Sprintf(url_format, b.APIKey, term))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	var data synonyms
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if data.Noun != nil {
		syn = append(syn, data.Noun.Syn...)
	}
	return syn, nil
}
