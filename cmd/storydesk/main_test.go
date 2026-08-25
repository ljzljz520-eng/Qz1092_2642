package main

import (
	"net/http/httptest"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	response := httptest.NewRecorder()
	healthHandler(response, httptest.NewRequest("GET", "/health", nil))
	if response.Code != 200 || response.Body.String() == "" {
		t.Fatalf("response=%d", response.Code)
	}
}
