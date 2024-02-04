package server

import (
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	port int
}

// New TODO: implement request logging (similar to middleware in Chi)
func New(port int) *http.Server {
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
