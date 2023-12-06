package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"auth/internal/app"
	"auth/internal/database"
	"auth/internal/logger"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.IndexHandler)
	mux.HandleFunc("/api/data", s.ApiDataHandler)
	mux.HandleFunc("/api/auth", s.LoginHandler)

	return mux
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cred app.Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	// Validate body and resulting struct
	if err == io.EOF || cred.Email == "" || cred.Password == "" {
		http.Error(w, "Invalid body", http.StatusBadGateway)

		return
	}

	logger.Debug("Credentials:", cred)

	c := &app.LoginController{
		UserRepository: &database.UserRepositoryImpl{},
	}
	response, err := c.Execute(cred)
	// Handle custom UserNotFoundError with specific response
	if err != nil {
		// if _, ok := err.(*app.UserNotFoundError); ok {
		// 	http.Error(w, "User not found", http.StatusNotFound)
		// 	return
		// }
		// http.Error(w, "Bad gateway", http.StatusBadGateway)
		// return

		logger.Error("Error:", err)

		switch err.(type) {
		case *app.UserNotFoundError:
			http.Error(w, "User not found", http.StatusNotFound)

			return
		default:
			http.Error(w, "Bad gateway", http.StatusBadGateway)

			return
		}
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["message"] = "Hello World"

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (s *Server) ApiDataHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"data": "some important data",
	}

	fmt.Println(response)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
