package weather

import "context"

type WeatherData struct {
	TempMax          []float64
	TempMin          []float64
	Pricipitation    []float64
	SunshineDuration []float64
}

type Coordinates struct {
	Longitude float64
	Latitude  float64
}

type Client interface {
	FetchCurrentWeatherData(ctx context.Context, coords []Coordinates) (*WeatherData, error)
}
