package main

import (
	"auth/internal/config"
	_ "auth/internal/config"
	"auth/internal/db"
	_ "auth/internal/db"
	"auth/internal/logger"
	"auth/internal/server"
	"fmt"
)

func main() {
	config.Init()
	database := db.Connect()
	s := server.New(config.Values.Port, database)

	logger.Info.Printf("Server listening on 0.0.0.0%s", s.Addr)

	err := s.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start")

		panic(err)
	}
}
