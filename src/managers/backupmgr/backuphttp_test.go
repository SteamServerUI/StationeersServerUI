package backupmgr

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"
)

func TestBackupHTTPSelectionUsesName(t *testing.T) {
	path := analysisFixture(t)
	m := NewBackupManager(BackupConfig{SafeBackupDir: filepath.Dir(path)})
	handler := &HTTPHandler{manager: m}
	query := url.Values{"name": {filepath.Base(path)}}
	request := httptest.NewRequest(http.MethodGet, "/api/v3/backups/analysis?"+query.Encode(), nil)
	response := httptest.NewRecorder()
	handler.AnalyzeBackupHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("stable selection depended on the old index: %d %s", response.Code, response.Body.String())
	}
	response = httptest.NewRecorder()
	handler.DownloadBackupHandler(response, httptest.NewRequest(http.MethodGet, "/api/v3/backups/download?"+query.Encode(), nil))
	if response.Code != http.StatusOK || !bytes.HasPrefix(response.Body.Bytes(), []byte("PK")) {
		t.Fatalf("stable download: %d", response.Code)
	}
	if _, err := m.GetBackupFileData("../outside.save"); err == nil {
		t.Fatal("accepted a file outside the inventory")
	}
}

func TestAnalyzeBackupHandler(t *testing.T) {
	path := analysisFixture(t)
	manager := NewBackupManager(BackupConfig{SafeBackupDir: filepath.Dir(path)})
	handler := &HTTPHandler{manager: manager}
	request := httptest.NewRequest(http.MethodGet, "/api/v3/backups/analysis?name=analysis.save", nil)
	response := httptest.NewRecorder()

	handler.AnalyzeBackupHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
	var body struct {
		Data SaveAnalysis `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	analysis := body.Data
	if analysis.DaysPlayed != 67 || analysis.Players != 3 || analysis.Furnaces != 2 {
		t.Fatalf("unexpected analysis response: %+v", analysis)
	}
}

func TestAnalyzeBackupHandlerRejectsInvalidIndex(t *testing.T) {
	handler := &HTTPHandler{manager: NewBackupManager(BackupConfig{})}
	request := httptest.NewRequest(http.MethodGet, "/api/v3/backups/analysis?index=nope", nil)
	response := httptest.NewRecorder()

	handler.AnalyzeBackupHandler(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status %d: %s", response.Code, response.Body.String())
	}
}
