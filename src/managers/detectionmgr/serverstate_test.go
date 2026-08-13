package detectionmgr

import "testing"

func TestLifecycleLogLinesAreDetectedIndividually(t *testing.T) {
	detector := NewDetector()
	var detected []EventType
	SetServerStateHandler(func(eventType EventType) {
		detected = append(detected, eventType)
	})
	t.Cleanup(func() { SetServerStateHandler(nil) })

	lines := []string{
		"00:31:03: game manager initialized",
		"00:31:56: RakNet successfully hosted with Address: 0.0.0.0:27016",
		"00:31:56: StartSession. config: {",
		"00:31:57: registered with session #7668",
	}
	for _, line := range lines {
		detector.ProcessLogMessage(line)
	}

	want := []EventType{EventGameManagerReady, EventServerHosted, EventSessionStarting, EventSessionRegistered}
	if len(detected) != len(want) {
		t.Fatalf("expected %d lifecycle detections, got %d: %v", len(want), len(detected), detected)
	}
	for i := range want {
		if detected[i] != want[i] {
			t.Fatalf("detection %d: expected %s, got %s", i, want[i], detected[i])
		}
	}
}
