package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

type Location struct {
	City      string  `json:"city"`
	ISO3      string  `json:"iso3"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading .env file.", err)
		return nil, err
	}

	config := &Config{
		Port:       os.Getenv("PORT"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
	}

	return config, nil
}

func LoadLocations(filePath string) ([]Location, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := json.Unmarshal(data, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}
