package worker

import (
	"context"
	"fmt"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/azlan-code/vibes-anxiety/integrations/weather"
	"github.com/azlan-code/vibes-anxiety/internal/scorer"
)

func prepareScoringData(ctx context.Context, weatherClient weather.Client, pastDays int) ([]scorer.ScoringData, error) {
	locations, err := config.LoadLocations("config/locations.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load locations: %w", err)
	}

	coords := make([]weather.Coordinates, len(locations))
	for i, loc := range locations {
		coords[i] = weather.Coordinates{
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
		}
	}

	wd, err := weatherClient.FetchWeatherData(ctx, coords, pastDays)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}

	var allScoringData []scorer.ScoringData
	for i, loc := range locations {
		allScoringData = append(allScoringData, scorer.ScoringData{
			Country: loc.ISO3,
			City:    loc.City,
			Coordinates: weather.Coordinates{
				Longitude: loc.Longitude,
				Latitude:  loc.Latitude,
			},
			Weather: wd[i],
		})
	}

	return allScoringData, nil
}

func CalculateVibeScores(ctx context.Context, pastDays int) ([]scorer.ScoreResult, error) {
	weatherClient := weather.NewHTTPClient()
	allScoringData, err := prepareScoringData(ctx, weatherClient, pastDays)
	if err != nil {
		return nil, err
	}

	vibeScores := scorer.CalculateVibeScores(allScoringData, pastDays)
	// s, _ := json.MarshalIndent(vibeScores, "", "  ")
	// fmt.Print(string(s))

	return vibeScores, nil
}
