package ssestream

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSSEPayloadIsJSON(t *testing.T) {
	response := httptest.NewRecorder()
	if err := writeSSEEvent(response, "message", "hello"); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(response.Body.String(), "\n")
	var data string
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			data = strings.TrimPrefix(line, "data: ")
		}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		t.Fatalf("SSE data is not JSON: %v", err)
	}
	if payload["type"] != "message" || payload["data"] != "hello" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}
