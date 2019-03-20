package main

import (
	"log"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatalln(err)
	}
}
