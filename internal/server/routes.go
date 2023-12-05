package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var baseRoute string = "/api/"

func makePattern(p string) string {
	return baseRoute + p
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.IndexHandler)
	mux.HandleFunc(makePattern("api/data"), s.ApiDataHandler)

	return mux
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
	response := map[string]string{}

	fmt.Println(response)

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
