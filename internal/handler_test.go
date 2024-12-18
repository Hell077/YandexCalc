package internal

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           interface{}
		expectedCode   int
		expectedResult *float64
		expectedError  *string
	}{
		{
			name:           "Valid Expression",
			method:         http.MethodPost,
			body:           map[string]string{"expression": "2+2"},
			expectedCode:   http.StatusOK,
			expectedResult: floatPointer(4),
		},
		{
			name:          "Invalid JSON",
			method:        http.MethodPost,
			body:          "invalid_json",
			expectedCode:  http.StatusBadRequest,
			expectedError: stringPointer("Invalid JSON format"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/api/v1/calculate", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			CalculateHandler(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("expected status %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			var responseBody Response
			_ = json.NewDecoder(resp.Body).Decode(&responseBody)

			if tt.expectedResult != nil && (responseBody.Result == nil || *responseBody.Result != *tt.expectedResult) {
				t.Errorf("expected result %v, got %v", tt.expectedResult, responseBody.Result)
			}

			if tt.expectedError != nil && (responseBody.Error == nil || *responseBody.Error != *tt.expectedError) {
				t.Errorf("expected error %s, got %v", *tt.expectedError, responseBody.Error)
			}
		})
	}
}

func floatPointer(f float64) *float64 {
	return &f
}

func stringPointer(s string) *string {
	return &s
}
