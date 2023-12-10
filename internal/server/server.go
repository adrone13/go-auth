package server

import (
	"auth/internal/config"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	port int
}

// TODO: implement request logging (similar to middleware in Chi)
func NewServer() *http.Server {
	// TODO: implement Config service to manage env variables
	port := config.GetInt("PORT")

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
