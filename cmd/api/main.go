package main

import (
	"context"
	"log"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/azlan-code/vibes-anxiety/internal/server"
	"github.com/azlan-code/vibes-anxiety/internal/worker"
)

func main() {
	ctx := context.Background()
	worker.CalculateVibeScore(ctx)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration")
	}

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("failed to create new server.")
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start to server.")
	}
}
