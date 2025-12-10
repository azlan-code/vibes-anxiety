package database

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	GetDB() *pgxpool.Pool
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

func (s *service) GetDB() *pgxpool.Pool {
	return s.db
}

func (s *service) Close() {
	log.Printf("Disconnected from database.")
	s.db.Close()
}
