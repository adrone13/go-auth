package main

import (
	"auth/internal/config"
	"auth/internal/logger"
	"auth/internal/server"
	"fmt"
)

func main() {
	s := server.New(config.Values.Port)

	logger.Info(
		fmt.Sprintf("Server listening on 0.0.0.0%s", s.Addr),
	)

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start")

		panic(err)
	}
}
