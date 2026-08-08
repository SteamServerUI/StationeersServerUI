package runfileapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/logger"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/runfile"
)

const maxEditableFileBytes = 10 << 20

type FileRequest struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type FileInfo struct {
	Filename    string `json:"filename"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func GetFileList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fileError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET requests are allowed")
		return
	}
	files := runfile.GetFiles()
	if files == nil {
		fileError(w, http.StatusNotFound, "runfile_missing", "No runfile is loaded or it has no editable files")
		return
	}
	items := make([]FileInfo, 0, len(files))
	for _, file := range files {
		if protectedFile(file.Filepath) {
			continue
		}
		if _, err := os.Stat(file.Filepath); err != nil {
			continue
		}
		items = append(items, FileInfo{Filename: file.Filename, Type: file.Type, Description: file.Description})
	}
	writeFileJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"files": items}})
}

func HandleFileContent(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if !hasFilePermission(r, security.PermissionFilesRead) {
			fileError(w, http.StatusForbidden, "forbidden", "Permission denied")
			return
		}
		readFile(w, r.URL.Query().Get("filename"))
	case http.MethodPut:
		if !hasFilePermission(r, security.PermissionFilesWrite) {
			fileError(w, http.StatusForbidden, "forbidden", "Permission denied")
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxEditableFileBytes)
		var request FileRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&request); err != nil {
			fileError(w, http.StatusBadRequest, "invalid_json", "A filename and textual content are required")
			return
		}
		saveFile(w, request)
	default:
		fileError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Only GET and PUT requests are allowed")
	}
}

func readFile(w http.ResponseWriter, filename string) {
	target, ok := editableFile(filename)
	if !ok {
		fileError(w, http.StatusNotFound, "file_not_found", "The requested file is not available in this runfile")
		return
	}
	reader, err := os.Open(target.Filepath)
	if err != nil {
		fileError(w, http.StatusNotFound, "file_unavailable", "The requested file could not be opened")
		return
	}
	defer reader.Close()
	content, err := io.ReadAll(io.LimitReader(reader, maxEditableFileBytes+1))
	if err != nil {
		fileError(w, http.StatusInternalServerError, "file_read_failed", "The requested file could not be read")
		return
	}
	if len(content) > maxEditableFileBytes {
		fileError(w, http.StatusRequestEntityTooLarge, "file_too_large", "The file is too large for the web editor")
		return
	}
	writeFileJSON(w, http.StatusOK, map[string]any{"data": map[string]any{
		"filename": target.Filename,
		"type":     target.Type,
		"content":  string(content),
	}})
}

func saveFile(w http.ResponseWriter, request FileRequest) {
	target, ok := editableFile(request.Filename)
	if !ok {
		fileError(w, http.StatusNotFound, "file_not_found", "The requested file is not available in this runfile")
		return
	}
	if len(request.Content) > maxEditableFileBytes {
		fileError(w, http.StatusRequestEntityTooLarge, "file_too_large", "The file is too large for the web editor")
		return
	}
	if info, err := os.Stat(target.Filepath); err == nil && info.Mode().Perm()&0222 == 0 {
		fileError(w, http.StatusForbidden, "file_read_only", "The requested file is read-only")
		return
	} else if err != nil && !os.IsNotExist(err) {
		fileError(w, http.StatusInternalServerError, "file_stat_failed", "The requested file could not be inspected")
		return
	}
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		err = os.WriteFile(target.Filepath, []byte(request.Content), 0644)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		logger.Runfile.Warn(fmt.Sprintf("failed to write editable file %s: %v", request.Filename, err))
		fileError(w, http.StatusInternalServerError, "file_save_failed", "The file could not be saved")
		return
	}
	writeFileJSON(w, http.StatusOK, map[string]any{"data": map[string]any{"filename": target.Filename, "saved": true}})
}

func editableFile(filename string) (*runfile.File, bool) {
	if strings.TrimSpace(filename) == "" {
		return nil, false
	}
	for _, file := range runfile.GetFiles() {
		if file.Filename == filename && !protectedFile(file.Filepath) {
			copy := file
			return &copy, true
		}
	}
	return nil, false
}

func protectedFile(path string) bool {
	normalized := strings.TrimPrefix(strings.ReplaceAll(path, "\\", "/"), "./")
	return normalized == "SSUI" || strings.HasPrefix(normalized, "SSUI/")
}

func hasFilePermission(r *http.Request, permission string) bool {
	principal, ok := middleware.PrincipalFromContext(r.Context())
	return ok && principal.Permissions[permission]
}

func fileError(w http.ResponseWriter, status int, code, message string) {
	logger.API.Debug(message)
	writeFileJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func writeFileJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
