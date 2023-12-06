package main

import (
	"fmt"

	"auth/internal/logger"
	"auth/internal/server"
)

func main() {
	server := server.NewServer()

	logger.Info(
		fmt.Sprintf("Server listening on 0.0.0.0%s", server.Addr),
	)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start")

		panic(err)
	}
}
