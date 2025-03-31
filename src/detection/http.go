package detection

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

/*
HTTP API for custom detections.
- Handles GET (list), POST (add), and DELETE (remove) requests for custom patterns.
- Interfaces with the CustomDetectionsManager over HTTP.
*/

var customDetectionsManager *CustomDetectionsManager

// InitCustomDetectionsManager initializes the custom detections manager
func InitCustomDetectionsManager(detector *Detector) {
	customDetectionsManager = NewCustomDetectionsManager(detector)
}

// HandleCustomDetections handles GET and POST requests for custom detections
func HandleCustomDetection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// Return all custom detections
		detections := customDetectionsManager.GetDetections()
		json.NewEncoder(w).Encode(detections)

	case http.MethodPost:
		// Add a new custom detection
		var newDetection CustomDetection
		if err := json.NewDecoder(r.Body).Decode(&newDetection); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Generate UUID if not provided
		if newDetection.ID == "" {
			newDetection.ID = uuid.New().String()
		}

		// Validate the detection
		if newDetection.Type != "regex" && newDetection.Type != "keyword" {
			http.Error(w, "Type must be 'regex' or 'keyword'", http.StatusBadRequest)
			return
		}

		if newDetection.Pattern == "" {
			http.Error(w, "Pattern cannot be empty", http.StatusBadRequest)
			return
		}

		// Default to CUSTOM_DETECTION if not specified
		if newDetection.EventType == "" {
			newDetection.EventType = "CUSTOM_DETECTION"
		}

		// Add the detection
		if err := customDetectionsManager.AddDetection(newDetection); err != nil {
			if strings.Contains(err.Error(), "Invalid Detection") {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		// Return the created detection
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newDetection)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleCustomDetectionsWithID handles DELETE requests for a specific custom detection
func HandleDeleteCustomDetection(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from query parameter
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing detection ID in query parameter (?id=)", http.StatusBadRequest)
		return
	}

	if id == "" {
		http.Error(w, "Missing detection ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		// Delete the detection
		if err := customDetectionsManager.DeleteDetection(id); err != nil {
			if strings.Contains(err.Error(), "detection not found") {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
