package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type VibeRecord struct {
	Longitude float64 `json:"longitude"` // Longitude coordinate
	Latitude  float64 `json:"latitude"`  // Latitude coordinate
	ISO3      string  `json:"iso3"`      // Country ISO3 code
	City      string  `json:"city"`      // City name
	Day       string  `json:"day"`       // YYYY-MM-DD
	Score     float64 `json:"score"`     // Vibe score
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/vibes/timeline", s.TimelineHandler)
	r.Post("/api/vibes", s.VibePostHandler)

	return r
}

func (s *Server) TimelineHandler(w http.ResponseWriter, r *http.Request) {
	db := s.db.GetDB()
	ctx := r.Context()

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	query := `
        SELECT longitude, latitude, iso3, city, day, score
        FROM location_vibes
        WHERE day >= $1
        ORDER BY iso3, day
    `

	rows, err := db.Query(ctx, query, thirtyDaysAgo)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []VibeRecord
	for rows.Next() {
		var record VibeRecord
		var dayTime time.Time
		err := rows.Scan(&record.Longitude, &record.Latitude, &record.ISO3, &record.City, &dayTime, &record.Score)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		record.Day = dayTime.Format("2006-01-02")
		records = append(records, record)
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(records)
	if err != nil {
		http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResp)
}

func (s *Server) VibePostHandler(w http.ResponseWriter, r *http.Request) {
	db := s.db.GetDB()
	ctx := r.Context()

	var vibe VibeRecord
	err := json.NewDecoder(r.Body).Decode(&vibe)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if vibe.Latitude == 0 && vibe.Longitude == 0 {
		http.Error(w, "latitude and longitude are required", http.StatusBadRequest)
		return
	}

	if vibe.ISO3 == "" || vibe.Day == "" || vibe.City == "" {
		http.Error(w, "iso3, city, and day are required", http.StatusBadRequest)
		return
	}

	// Parse the day string into time.Time
	dayTime, err := time.Parse("2006-01-02", vibe.Day)
	if err != nil {
		http.Error(w, "Invalid day format, expected YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO location_vibes (iso3, city, longitude, latitude, day, score)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (longitude, latitude, day) DO UPDATE SET score = EXCLUDED.score
    `

	_, err = db.Exec(ctx, query, vibe.ISO3, vibe.City, vibe.Longitude, vibe.Latitude, dayTime, vibe.Score)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
