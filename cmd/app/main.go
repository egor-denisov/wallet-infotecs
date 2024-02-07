package main

import (
	"log"

	"github.com/egor-denisov/wallet-infotecs/config"
	"github.com/egor-denisov/wallet-infotecs/internal/app"
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
