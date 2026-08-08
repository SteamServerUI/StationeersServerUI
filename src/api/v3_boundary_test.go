package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestV3BoundaryNormalizesPlainResponses(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		body       string
		wantStatus int
		wantKey    string
	}{
		{name: "success", status: http.StatusOK, body: "started", wantStatus: http.StatusOK, wantKey: "data"},
		{name: "failure", status: http.StatusBadRequest, body: "bad request", wantStatus: http.StatusBadRequest, wantKey: "error"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := v3JSONBoundary(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(test.status)
				_, _ = w.Write([]byte(test.body))
			})
			response := httptest.NewRecorder()
			handler(response, httptest.NewRequest(http.MethodGet, "/api/v3/example", nil))
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Fatalf("content type = %q", response.Header().Get("Content-Type"))
			}
			var payload map[string]any
			if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
				t.Fatalf("response is not JSON: %v", err)
			}
			if payload[test.wantKey] == nil {
				t.Fatalf("response has no %q envelope: %#v", test.wantKey, payload)
			}
		})
	}
}

func TestV3BoundaryPreservesJSON(t *testing.T) {
	handler := v3JSONBoundary(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]bool{"ready": true}})
	})
	response := httptest.NewRecorder()
	handler(response, httptest.NewRequest(http.MethodGet, "/api/v3/example", nil))
	var payload struct {
		Data struct {
			Ready bool `json:"ready"`
		} `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil || !payload.Data.Ready {
		t.Fatalf("JSON payload was not preserved: %s (%v)", response.Body.String(), err)
	}
}
