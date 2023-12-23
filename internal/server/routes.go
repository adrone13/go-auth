package server

import (
	"encoding/json"
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
	response["server"] = "running 🚀"

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
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	u := app.GetUserUseCase{
		UserRepository: &db.UserRepositoryImpl{},
	}
	user := u.Execute(app.GetUserInput{
		Id: app.UserId(token.Claims.Subject),
	})

	jsonResp, err := json.Marshal(user)
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
	// Validate input struct
	if err == io.EOF || cred.Email == "" || cred.Password == "" {
		http.Error(w, "Invalid body", http.StatusBadGateway)

		return
	}

	logger.Debug("Credentials:", cred)

	c := &app.LogInUseCase{
		UserRepository: &db.UserRepositoryImpl{},
	}
	response, err := c.Execute(cred)

	if err != nil {
		switch err.(type) {
		case *app.UserNotFoundError:
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
