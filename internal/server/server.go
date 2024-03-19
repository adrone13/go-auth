package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Port int
	DB   Database
}

// New TODO: implement request logging (similar to middleware in Chi)
func New(port int, db Database) *http.Server {
	NewServer := &Server{
		Port: port,
		DB:   db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
