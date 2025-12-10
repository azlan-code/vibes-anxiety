package scorer

import "github.com/azlan-code/vibes-anxiety/integrations/weather"

type ScoringData struct {
	Country     string
	City        string
	Coordinates weather.Coordinates
	Weather     weather.WeatherData
	// Add other integration data sources here
	// GPR []gpr.GPRData
}

type ScoreResult struct {
	Country     string
	City        string
	Coordinates weather.Coordinates
	Scores      []float64
}

func CalculateVibeScores(data []ScoringData, pastDays int) []ScoreResult {
	var results []ScoreResult
	for _, d := range data {
		var scores []float64
		for i := range pastDays {
			tempMax := getValue(d.Weather.TempMax, i)
			tempMin := getValue(d.Weather.TempMin, i)
			precipitation := getValue(d.Weather.Pricipitation, i)
			sunshineDuration := getValue(d.Weather.SunshineDuration, i)

			tempMaxScore := calculateVibeForAttribute(25, 10, 32, 26, tempMax)
			tempMinScore := calculateVibeForAttribute(25, 0, 28, 18, tempMin)
			precipitationScore := calculateVibeForAttribute(25, 0, 2, 0, precipitation)
			sunshineDurationScore := calculateVibeForAttribute(25, 0, 86400, 43200, sunshineDuration)

			totalScore := sum(
				tempMaxScore,
				tempMinScore,
				precipitationScore,
				sunshineDurationScore,
			)
			scores = append(scores, totalScore)
		}
		results = append(results, ScoreResult{
			Country:     d.Country,
			City:        d.City,
			Coordinates: d.Coordinates,
			Scores:      scores,
		})
	}

	// Combines weather data, gpr data, etc.
	return results
}

func getValue(slice []float64, i int) float64 {
	if i < len(slice) {
		return slice[i]
	}
	return 0
}

func calculateVibeForAttribute(total float64, lowerLimit float64, upperLimit float64, perfect float64, value float64) float64 {
	if value < lowerLimit || value > upperLimit {
		return 0
	}
	if value == perfect {
		return total
	}
	if value < perfect {
		return total * (value - lowerLimit) / (perfect - lowerLimit)
	}
	return total * (upperLimit - value) / (upperLimit - perfect)
}

func sum(vals ...float64) float64 {
	total := 0.0
	for _, v := range vals {
		total += v
	}
	return total
}
