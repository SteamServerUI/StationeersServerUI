package detection

// This file contains the public API for the detection package

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
