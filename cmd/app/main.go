package main

import (
	"log"

	"github.com/ev-go/Testing/config"
	"github.com/ev-go/Testing/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
