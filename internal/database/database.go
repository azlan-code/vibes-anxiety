package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetDB() *pgxpool.Pool
	AddVibeRecord(ctx context.Context, iso3 string, city string, longitude float64, latitude float64, day time.Time, score float64) error
	IsEmpty(ctx context.Context) (bool, error)
	PopulateScores(ctx context.Context, iso3 string, city string, longitude, latitude float64, scores []float64, startDate time.Time) error
	Close()
}

type service struct {
	db *pgxpool.Pool
}

//go:embed migrations/001_create_schema.sql
var migrationSQL string

var dbInstance *service

func NewDatabase(cfg *config.Config) Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=public",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSLMode,
	)
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to create new db pool: %v", err)
	}

	ctx := context.Background()
	err = ensureTimescale(ctx, pool)
	if err != nil {
		log.Fatalf("timescale migration failed: %v", err)
	}

	dbInstance = &service{
		db: pool,
	}
	return dbInstance
}

func ensureTimescale(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, migrationSQL)
	return err
}

func (s *service) AddVibeRecord(ctx context.Context, iso3 string, city string, longitude float64, latitude float64, day time.Time, score float64) error {
	query := `
		INSERT INTO location_vibes (iso3, city, longitude, latitude, day, score)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (longitude, latitude, day) DO UPDATE SET score = EXCLUDED.score
	`

	_, err := s.db.Exec(ctx, query, iso3, city, longitude, latitude, day, score)
	if err != nil {
		return fmt.Errorf("failed to add vibe record: %w", err)
	}

	return nil
}

func (s *service) PopulateScores(ctx context.Context, iso3 string, city string, longitude, latitude float64, scores []float64, startDate time.Time) error {
	for i, score := range scores {
		day := startDate.AddDate(0, 0, i)
		err := s.AddVibeRecord(ctx, iso3, city, longitude, latitude, day, score)
		if err != nil {
			return fmt.Errorf("failed to add score for day %d: %w", i, err)
		}
	}
	return nil
}

func (s *service) IsEmpty(ctx context.Context) (bool, error) {
	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM location_vibes").Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if database is empty: %w", err)
	}
	return count == 0, nil
}

func (s *service) GetDB() *pgxpool.Pool {
	return s.db
}

func (s *service) Close() {
	log.Printf("Disconnected from database.")
	s.db.Close()
}
