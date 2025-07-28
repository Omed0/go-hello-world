package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/omed0/go-hello-world/handlers"
)

// TestHandlerReadiness tests the health check endpoint
func TestHandlerReadiness(t *testing.T) {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandlerReadiness)

	// Call the handler with our request and recorder
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body contains the expected JSON
	expected := `{"status":"healthy","message":"Service is running properly"}`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestHandlerErr tests the error endpoint
func TestHandlerErr(t *testing.T) {
	req, err := http.NewRequest("GET", "/err", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HandlerErr)

	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	// Check that we get an error response
	if rr.Body.String() == "" {
		t.Error("Handler returned empty body")
	}
}

// Example of how to run tests:
// go test ./handlers -v
