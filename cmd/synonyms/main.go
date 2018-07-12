package main

import (
	"os"
	"log"
	"chatter/thesaurus"
	"bufio"
	"fmt"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	if apiKey == "" {
		log.Fatalln("error: empty epi key")
	}
	t := thesaurus.New(apiKey)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := t.Synonyms(word)
		if err != nil {
			log.Fatal(err.Error())
		}
		if len(syns) == 0 {
			log.Fatalf("synonyms: couldn't find any synonym for word '%s'", word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
