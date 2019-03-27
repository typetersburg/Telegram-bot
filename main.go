package main

import (
	"log"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatalln(err)
	}

	bot, err := cfg.Tg.New()
	if err != nil {
		log.Fatalf("Cannot initialize connection to tg bot api: %s", err)
	}
}
