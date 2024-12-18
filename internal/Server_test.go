package internal

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	address := "localhost:8080"
	server := NewServer(address)

	if server.address != address {
		t.Errorf("expected address %q, got %q", address, server.address)
	}
}

func TestServerRun(t *testing.T) {
	_ = NewServer("localhost:8080")

	testServer := httptest.NewServer(http.HandlerFunc(CalculateHandler))
	defer testServer.Close()

	resp, err := http.Post(testServer.URL+"/api/v1/calculate", "application/json", nil)
	if err != nil {
		t.Fatalf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	expectedError := "Invalid JSON format"
	if response["error"] != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, response["error"])
	}
}
