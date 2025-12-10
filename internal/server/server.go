package server

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/azlan-code/vibes-anxiety/config"
	"github.com/azlan-code/vibes-anxiety/internal/database"
)

type Server struct {
	port int
	db   database.Service
}

func NewServer(cfg *config.Config) (*http.Server, error) {
	port, err := strconv.Atoi(cfg.Port)
	if err != nil {
		log.Println("failed to get PORT from env.", err)
		return nil, err
	}

	NewServer := &Server{
		port: port,
		db:   database.NewDatabase(cfg),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", NewServer.port),
		Handler: NewServer.RegisterRoutes(),
	}

	return server, nil
}
