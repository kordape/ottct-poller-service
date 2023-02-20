package main

import (
	"log"

	"github.com/kordape/ottct-poller-service/config"
)

func main() {
	// Configuration
	_, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
}
