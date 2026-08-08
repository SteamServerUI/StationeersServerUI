package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPluginRoutesAreQuarantinedByDefault(t *testing.T) {
	t.Setenv("SSUI_ENABLE_UNSAFE_PLUGINS", "")
	_, protected := SetupV7APIRoutes()
	response := httptest.NewRecorder()
	protected.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/v3/plugingallery/select", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("plugin endpoint status = %d, want %d", response.Code, http.StatusServiceUnavailable)
	}
}

func TestStaticFilesRequireFilesRead(t *testing.T) {
	t.Setenv("SSUI_ENABLE_UNSAFE_PLUGINS", "")
	_, protected := SetupV7APIRoutes()
	response := httptest.NewRecorder()
	protected.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/files/example.txt", nil))
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("static file status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}
