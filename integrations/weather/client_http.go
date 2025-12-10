package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type httpClient struct {
	baseUrl string
	http    *http.Client
}

func NewHTTPClient() Client {
	return &httpClient{
		baseUrl: "https://api.open-meteo.com/v1/forecast",
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

type WeatherResponse struct {
	Daily struct {
		TempMax          []float64 `json:"temperature_2m_max"`
		TempMin          []float64 `json:"temperature_2m_min"`
		Pricipitation    []float64 `json:"precipitation_sum"`
		SunshineDuration []float64 `json:"sunshine_duration"`
	} `json:"daily"`
}

func (c *httpClient) FetchCurrentWeatherData(ctx context.Context, coords []Coordinates) (*WeatherData, error) {
	longitudes := make([]string, len(coords))
	latitudes := make([]string, len(coords))
	for i, c := range coords {
		longitudes[i] = strconv.FormatFloat(c.Longitude, 'f', -1, 64)
		latitudes[i] = strconv.FormatFloat(c.Latitude, 'f', -1, 64)
	}

	u, _ := url.Parse(c.baseUrl)
	q := u.Query()
	q.Set("longitude", strings.Join(longitudes, ","))
	q.Set("latitude", strings.Join(latitudes, ","))
	q.Set("daily", strings.Join([]string{
		"apparent_temperature_max",
		"apparent_temperature_min",
		"precipitation_sum",
		"sunshine_duration",
	}, ","))
	q.Set("past_days", "30")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	fmt.Println(u.String())

	var decoded WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return nil, err
	}

	return &WeatherData{
		TempMax:          decoded.Daily.TempMax,
		TempMin:          decoded.Daily.TempMin,
		Pricipitation:    decoded.Daily.Pricipitation,
		SunshineDuration: decoded.Daily.SunshineDuration,
	}, nil
}
