package setup

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func TestDownloadFileWithClient(t *testing.T) {
	target := filepath.Join(t.TempDir(), "downloads", "bepinex.zip")
	client := &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Status:     "200 OK",
				Body:       io.NopCloser(strings.NewReader("bepinex archive")),
			}, nil
		}),
	}
	if err := downloadFileWithClient(client, target, "https://example.invalid/bepinex.zip"); err != nil {
		t.Fatalf("downloadFileWithClient returned an error: %v", err)
	}

	contents, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read downloaded file: %v", err)
	}
	if string(contents) != "bepinex archive" {
		t.Fatalf("unexpected downloaded contents: %q", contents)
	}
}

func TestDownloadFileWithClientRejectsBadStatus(t *testing.T) {
	target := filepath.Join(t.TempDir(), "bepinex.zip")
	client := &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Status:     "503 Service Unavailable",
				Body:       io.NopCloser(strings.NewReader("unavailable")),
			}, nil
		}),
	}
	err := downloadFileWithClient(client, target, "https://example.invalid/bepinex.zip")
	if err == nil || !strings.Contains(err.Error(), "503 Service Unavailable") {
		t.Fatalf("expected status error, got %v", err)
	}
}
