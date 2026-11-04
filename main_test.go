package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMainRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	tests := []struct {
		name           string
		url            string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Health check",
			url:            "/health",
			expectedStatus: http.StatusOK,
			expectedBody:   `"status":"healthy"`,
		},
		{
			name:           "Employee exists",
			url:            "/employees/101",
			expectedStatus: http.StatusOK,
			expectedBody:   `"name":"Siva"`,
		},
		{
			name:           "Employee does not exist",
			url:            "/employees/999",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `"error":"Employee not found"`,
		},
		{
			name:           "Invalid employee ID",
			url:            "/employees/abc",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `"error":"Invalid employee ID"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(
				http.MethodGet,
				tt.url,
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf(
					"Expected status %d, got %d",
					tt.expectedStatus,
					recorder.Code,
				)
			}

			if !contains(recorder.Body.String(), tt.expectedBody) {
				t.Errorf(
					"Expected response to contain %q, got %s",
					tt.expectedBody,
					recorder.Body.String(),
				)
			}
		})
	}
}

func TestEmployeeResponseJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := setupRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/employees/101",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	var employee Employee

	err := json.Unmarshal(recorder.Body.Bytes(), &employee)

	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if employee.ID != 101 {
		t.Errorf("Expected employee ID 101, got %d", employee.ID)
	}

	if employee.Name != "Siva" {
		t.Errorf("Expected employee name Siva, got %s", employee.Name)
	}
}

// Small helper to check whether the response contains expected text.
func contains(body, expected string) bool {
	return len(body) >= len(expected) &&
		(body == expected ||
			len(body) > len(expected) &&
				search(body, expected))
}

func search(body, expected string) bool {
	for i := 0; i <= len(body)-len(expected); i++ {
		if body[i:i+len(expected)] == expected {
			return true
		}
	}

	return false
}
