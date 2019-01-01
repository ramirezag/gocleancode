package handlers_test

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"gocleancode/handlers"
	mockServices "gocleancode/services/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createHandlers() (*mockServices.FileService, *mux.Router) {
	fileService := &mockServices.FileService{}
	appHandlers := handlers.NewHandlers(fileService)
	return fileService, appHandlers
}

func TestStatus(t *testing.T) {
	_, appHandlers := createHandlers()
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// When
	appHandlers.ServeHTTP(rr, req)
	// Then
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expectedResponse := handlers.Response{true, "UP"}
	actualResponse := handlers.Response{}
	err = json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	if err != nil {
		t.Errorf("Expected no error, but got %s instead", err)
		return
	}
	assert.Equal(t, expectedResponse, actualResponse)
}