package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"auth/internal/app"
	"auth/internal/db"
	"auth/internal/logger"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(Get("/health", s.HealthHandler))
	mux.HandleFunc(Post("/api/login", s.LoginHandler))
	mux.HandleFunc(Post("/api/signup", s.SignUpHandler))

	return mux
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["server"] = "running 🚀"

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var i app.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&i)
	if err == io.EOF || i.Email == "" || i.FullName == "" || i.Password == "" {
		http.Error(w, "Invalid body", http.StatusBadGateway)

		return
	}

	useCase := &app.SignUpUseCase{
		UserRepository: &db.UserRepositoryImpl{},
	}

	err = useCase.Execute(i)
	if err != nil {
		logger.Error(err)

		http.Error(w, "Bad gateway", http.StatusBadGateway)

		return
	}

	w.WriteHeader(http.StatusOK)
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

	c := &app.LogInUseCase{
		UserRepository: &db.UserRepositoryImpl{},
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
