package main

import (
	"log"

	"github.com/mirjalilova/ccenter_news.git/config"
	"github.com/mirjalilova/ccenter_news.git/internal/app"
)

func main() {

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
