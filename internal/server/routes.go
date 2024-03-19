package server

import (
	"encoding/json"
	"errors"
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
	mux.HandleFunc(Get("/api/me", s.MeHandler))
	mux.HandleFunc(Post("/api/login", s.LoginHandler))
	mux.HandleFunc(Post("/api/signup", s.SignUpHandler))

	return mux
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)

	var dbStatus string
	if err := s.DB.Ping(r.Context()); err == nil {
		dbStatus = "running 🚀"
	} else {
		dbStatus = fmt.Sprintf("failing. error: %s", err)
	}

	response["server"] = "running 🚀"
	response["db"] = dbStatus

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func (s *Server) MeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := authenticate(w, r)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	u := app.GetUserUseCase{
		UserRepository: &db.UserRepository{},
	}
	user, err := u.Execute(r.Context(), app.UserId(token.Claims.Subject))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	jsonResp, _ := json.Marshal(user)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

// An example successful response:
//
// HTTP/1.1 200 OK
// Content-Type: application/json;charset=UTF-8
// Cache-Control: no-store
// Pragma: no-cache
//
// {
//     "access_token":"2YotnFZFEjr1zCsicMWpAA",
//     "token_type":"example",
//     "expires_in":3600,
//     "refresh_token":"tGzv3JOkF0XG5Qx2TlKWIA",
//     "example_parameter":"example_value"
// }

func (s *Server) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var i app.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&i)
	if err == io.EOF || i.Email == "" || i.FullName == "" || i.Password == "" {
		http.Error(w, "Invalid body", http.StatusBadGateway)

		return
	}

	useCase := &app.SignUpUseCase{
		UserRepo: &db.UserRepository{},
	}

	err = useCase.Execute(r.Context(), i)
	if err != nil {
		logger.Error(err)
		http.Error(w, "Bad gateway", http.StatusBadGateway)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var cred app.Credentials
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err == io.EOF || cred.Email == "" || cred.Password == "" {
		http.Error(w, "Invalid body", http.StatusBadGateway)

		return
	}

	logger.Debug(fmt.Sprintf("Credentials: %+v", cred))

	u := &app.LogInUseCase{
		UserRepository: &db.UserRepository{},
	}
	response, err := u.Execute(r.Context(), cred)

	if err != nil {
		logger.Error(err)

		var userNotFoundError *app.UserNotFoundError
		switch {
		case errors.As(err, &userNotFoundError):
			http.Error(w, "User not found", http.StatusNotFound)
		default:
			http.Error(w, "Bad gateway", http.StatusBadGateway)
		}

		return
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(jsonResp)
}
