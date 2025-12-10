package main

import (
	"context"
	"log"
	"time"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/azlan-code/vibes-anxiety/internal/database"
	"github.com/azlan-code/vibes-anxiety/internal/server"
	"github.com/azlan-code/vibes-anxiety/internal/worker"
)

func initializeDatabase(ctx context.Context, dbService database.Service) {
	// Check if database is empty and populate with initial data
	isEmpty, err := dbService.IsEmpty(ctx)
	if err != nil {
		log.Printf("warning: failed to check if database is empty: %v", err)
		return
	}

	if !isEmpty {
		return
	}

	log.Println("Database is empty, populating with initial data...")

	// Calculate vibe scores (30 days of past data)
	vibeScores, err := worker.CalculateVibeScores(ctx, 30)
	if err != nil {
		log.Printf("warning: failed to calculate initial vibe scores: %v", err)
		return
	}

	// Populate database with scores
	now := time.Now()
	for _, result := range vibeScores {
		// Calculate start date (30 days ago, or adjust based on your needs)
		startDate := now.AddDate(0, 0, -len(result.Scores)+1)

		err := dbService.PopulateScores(
			ctx,
			result.Country,
			result.City,
			result.Coordinates.Longitude,
			result.Coordinates.Latitude,
			result.Scores,
			startDate,
		)
		if err != nil {
			log.Printf("warning: failed to populate scores for %s: %v", result.City, err)
		}
	}

	log.Println("Initial data population completed")
}

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load configuration")
	}

	// Create database
	dbService := database.NewDatabase(cfg)
	defer dbService.Close()

	// Initialize database with initial data if empty
	initializeDatabase(ctx, dbService)

	server, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal("failed to create new server.")
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("failed to start to server.")
	}
}
