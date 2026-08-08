package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRemovedAPIRoutesReturnJSONNotFound(t *testing.T) {
	_, protected := SetupV7APIRoutes()
	response := httptest.NewRecorder()
	protected.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v2/server/status", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("old API status = %d, want %d", response.Code, http.StatusNotFound)
	}
	var payload map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("old API response is not JSON: %v", err)
	}
	if payload["error"] == nil {
		t.Fatalf("old API response has no error: %#v", payload)
	}
}
