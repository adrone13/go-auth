package server

import (
	"auth/internal/config"
	internalServer "auth/internal/server"
	"auth/tests/mocks"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	config.Values = mocks.ConfigMock

	s := internalServer.Server{Port: 8080, Db: &mocks.DatabaseMock{}}

	server := httptest.NewServer(http.HandlerFunc(s.HealthHandler))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	defer resp.Body.Close()

	// Assertions
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"server\":\"running 🚀\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}

	expectedHeader := "application/json"
	receivedHeader := resp.Header.Get("Content-Type")
	fmt.Println(resp.Header)
	if receivedHeader != expectedHeader {
		t.Errorf("expected \"Content-Type\" header to be %v; got %v", expectedHeader, receivedHeader)
	}
}
