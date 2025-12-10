package worker

import (
	"context"
	"fmt"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/azlan-code/vibes-anxiety/integrations/weather"
	"github.com/azlan-code/vibes-anxiety/internal/scorer"
)

func prepareScoringData(ctx context.Context, weatherClient weather.Client) (*scorer.ScoringData, error) {
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

	wd, err := weatherClient.FetchCurrentWeatherData(ctx, coords)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	fmt.Println(wd)

	scoringData := scorer.ScoringData{
		Weather: wd,
		// GPR: gprData,
		// ... other integration data
	}

	return &scoringData, nil
}

func CalculateVibeScore(ctx context.Context) error {
	weatherClient := weather.NewHTTPClient()
	scoringData, err := prepareScoringData(ctx, weatherClient)
	if err != nil {
		return err
	}

	vibeScore := scorer.CalculateVibeScore(*scoringData)
	fmt.Println(vibeScore)

	return nil
}
