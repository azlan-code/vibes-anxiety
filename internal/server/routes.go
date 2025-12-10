package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type VibeRecord struct {
	Longitude float64  `json:"longitude"` // Longitude coordinate
	Latitude  float64  `json:"latitude"`  // Latitude coordinate
	ISO3      string   `json:"iso3"`      // Country ISO3 code
	Day       string   `json:"day"`       // YYYY-MM-DD
	Score     *float64 `json:"score"`     // Vibe score
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

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	query := `
        SELECT coordinates[0] as longitude, coordinates[1] as latitude, iso3, day, score
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
		err := rows.Scan(&record.Longitude, &record.Latitude, &record.ISO3, &record.Day, &record.Score)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
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

	if vibe.ISO3 == "" || vibe.Day == "" {
		http.Error(w, "iso3 and day are required", http.StatusBadRequest)
		return
	}

	query := `
        INSERT INTO location_vibes (iso3, coordinates, day, score)
        VALUES ($1, POINT($2, $3), $4, $5)
        ON CONFLICT (coordinates, day) DO UPDATE SET score = EXCLUDED.score
    `

	_, err = db.Exec(ctx, query, vibe.Longitude, vibe.Latitude, vibe.ISO3, vibe.Day, vibe.Score)
	if err != nil {
		http.Error(w, "Error inserting data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
