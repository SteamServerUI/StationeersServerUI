package ssestream

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func writeSSEEvent(w http.ResponseWriter, eventType, data string) error {
	payload, err := json.Marshal(map[string]any{
		"type":       eventType,
		"data":       data,
		"occurredAt": time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "id: %d\nevent: message\ndata: %s\n\n", time.Now().UnixNano(), payload)
	return err
}
