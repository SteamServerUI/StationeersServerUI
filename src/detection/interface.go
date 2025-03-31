// interface.go
package detection

/*
Code-Public Detection API interface
- Exposes simplified interface for external references if needed
- Provides access to core detector functionality:
  - System initialization
  - New Handler registration - just in case we need to add more handlers later
  - Log processing
  - State queries (connected players currently)
*/

// Start initializes the detector and returns it
func Start() *Detector {
	return NewDetector()
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
