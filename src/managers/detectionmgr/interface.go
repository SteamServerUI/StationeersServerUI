// interface.go
package detectionmgr

import "sync"

/*
Code-Public Detection API interface
- Exposes simplified interface for external references if needed
- Provides access to core detector functionality:
  - System initialization
  - New Handler registration - just in case we need to add more handlers later
  - Log processing
  - State queries (connected players currently)
*/

var (
	detectorInstance *Detector
	once             sync.Once
)

// Start initializes the detector and stores it as the singleton instance
func Start() *Detector {
	once.Do(func() {
		detectorInstance = NewDetector()
	})
	return detectorInstance
}

// GetDetector returns the singleton detector instance
func GetDetector() *Detector {
	if detectorInstance == nil {
		panic("Detector not initialized. Call Start() first.")
	}
	return detectorInstance
}

// AddHandler is a convenient method to register a handler for an event type
func AddHandler(detector *Detector, eventType EventType, handler Handler) {
	detector.RegisterHandler(eventType, handler)
}

// ProcessLog is a convenient method to process a log message
func ProcessLog(detector *Detector, logMessage string) {
	detector.ProcessLogMessage(logMessage)
}

// GetPlayers returns the currently connected players
func GetPlayers(detector *Detector) map[string]string {
	return detector.GetConnectedPlayers()
}

// Clearplayers clears the connected players
func ClearPlayers(detector *Detector) {
	detector.ClearConnectedPlayers()
}
