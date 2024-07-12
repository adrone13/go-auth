package server

import (
	"auth/internal/app/users"
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
	mux.HandleFunc(Post("/api/token", s.RefreshToken))

	return mux
}

func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)

	var dbStatus string
	if err := s.Db.Ping(r.Context()); err == nil {
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
		logger.Error.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	u := app.GetUserUseCase{
		UserRepository: &db.UserRepository{},
	}
	user, err := u.Execute(r.Context(), users.UserId(token.Claims.Subject))
	if err != nil {
		http.NotFound(w, r)

		return
	}

	jsonResp, _ := json.Marshal(user)

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
		UserRepo: &db.UserRepository{},
	}

	err = useCase.Execute(r.Context(), i)
	if err != nil {
		logger.Error.Println(err)
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

	logger.Debug.Printf("Credentials: %+v", cred)

	u := &app.LogInUseCase{
		UserRepository: &db.UserRepository{},
	}
	response, err := u.Execute(r.Context(), cred)

	if err != nil {
		logger.Error.Println(err)

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

func (s *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	_, err := authenticate(w, r)
	if err != nil {
		logger.Error.Println(err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	grantType := r.URL.Query().Get("grant_type")
	refreshToken := r.URL.Query().Get("refresh_token")

	if grantType != "refresh_token" {
		http.Error(w, "invalid_request", http.StatusBadGateway)

		return
	}

	usecase := &app.RefreshAuthUseCase{
		UserRepository:     &db.UserRepository{},
		SessionsRepository: &db.SessionRepository{},
	}

	response, err := usecase.Execute(r.Context(), refreshToken)
	if err != nil {
		if err == errors.New("access_denied") {
			http.Error(w, "access_denied", http.StatusForbidden)

			return
		} else {
			http.Error(w, "invalid_request", http.StatusBadRequest)

			return
		}
	}

	jsonResp, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	w.Write(jsonResp)
}
