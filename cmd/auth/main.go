package main

import (
	"fmt"

	"auth/internal/server"
)

func main() {
	server := server.NewServer()

	fmt.Printf("Server listening on 0.0.0.0%s\n", server.Addr)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to start")

		panic(err)
	}
}
