package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/core/loader"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/steamserverui/runfile"
)

// APIGameArg is a DTO for GameArg, including RuntimeValue and all fields
type APIGameArg struct {
	Flag          string `json:"flag"`
	Value         string `json:"value"`
	RuntimeValue  string `json:"runtime_value"`
	Required      bool   `json:"required"`
	RequiresValue bool   `json:"requires_value"`
	Description   string `json:"description"`
	Type          string `json:"type"`
	Special       string `json:"special,omitempty"`
	UILabel       string `json:"ui_label"`
	UIGroup       string `json:"ui_group"`
	Weight        int    `json:"weight"`
	Min           int    `json:"min,omitempty"`
	Max           int    `json:"max,omitempty"`
	Disabled      bool   `json:"disabled"`
}

// APIMeta mirrors runfile.Meta for API responses
type APIMeta struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// APIRunFile is a DTO for RunFile, using APIGameArg and APIMeta
type APIRunFile struct {
	Meta              APIMeta                 `json:"meta"`
	Architecture      string                  `json:"architecture,omitempty"`
	SteamAppID        string                  `json:"steam_app_id"`
	WindowsExecutable string                  `json:"windows_executable"`
	LinuxExecutable   string                  `json:"linux_executable"`
	Args              map[string][]APIGameArg `json:"args"`
}

// apiResponse is the standard JSON response format
type apiResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error,omitempty"`
}

// toAPIGameArg converts runfile.GameArg to APIGameArg
func toAPIGameArg(arg runfile.GameArg) APIGameArg {
	return APIGameArg{
		Flag:          arg.Flag,
		Value:         arg.Value,
		RuntimeValue:  arg.RuntimeValue,
		Required:      arg.Required,
		RequiresValue: arg.RequiresValue,
		Description:   arg.Description,
		Type:          arg.Type,
		Special:       arg.Special,
		UILabel:       arg.UILabel,
		UIGroup:       arg.UIGroup,
		Weight:        arg.Weight,
		Min:           arg.Min,
		Max:           arg.Max,
		Disabled:      arg.Disabled,
	}
}

// toAPIRunFile converts runfile.RunFile to APIRunFile
func toAPIRunFile(rf *runfile.RunFile) APIRunFile {
	apiArgs := make(map[string][]APIGameArg)
	for category, args := range rf.Args {
		for _, arg := range args {
			apiArgs[category] = append(apiArgs[category], toAPIGameArg(arg))
		}
	}
	return APIRunFile{
		Meta: APIMeta{
			Name:    rf.Meta.Name,
			Version: rf.Meta.Version,
		},
		Architecture:      rf.Architecture,
		SteamAppID:        rf.SteamAppID,
		WindowsExecutable: rf.WindowsExecutable,
		LinuxExecutable:   rf.LinuxExecutable,
		Args:              apiArgs,
	}
}

// writeJSONResponse writes a JSON response with the given status code
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}, errMsg string) {
	resp := apiResponse{Data: data, Error: errMsg}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to write JSON response: %v", err))
	}
}

// HandleRunfileGroups handles GET /api/v2/runfile/groups
func HandleRunfileGroups(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("GET /api/v2/runfile/groups")
	if runfile.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	groups := runfile.GetUIGroups()
	logger.Runfile.Info("fetched UI groups")
	writeJSONResponse(w, http.StatusOK, groups, "")
}

// HandleRunfileArgs handles GET /api/v2/runfile/args
func HandleRunfileArgs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	logger.Runfile.Debug(fmt.Sprintf("GET /api/v2/runfile/args group=%s", group))

	if runfile.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	var args []runfile.GameArg
	if group != "" {
		// Validate group
		validGroups := runfile.GetUIGroups()
		valid := false
		for _, g := range validGroups {
			if g == group {
				valid = true
				break
			}
		}
		if !valid {
			logger.Runfile.Error(fmt.Sprintf("invalid group: %s", group))
			writeJSONResponse(w, http.StatusBadRequest, nil, fmt.Sprintf("invalid group: %s", group))
			return
		}
		args = runfile.GetArgsByGroup(group)
	} else {
		args = runfile.GetAllArgs()
	}

	// Convert to APIGameArg
	apiArgs := make([]APIGameArg, len(args))
	for i, arg := range args {
		apiArgs[i] = toAPIGameArg(arg)
	}

	logger.Runfile.Info(fmt.Sprintf("fetched args for group=%s", group))
	writeJSONResponse(w, http.StatusOK, apiArgs, "")
}

// HandleRunfileArgUpdate handles POST /api/v2/runfile/args/update
func HandleRunfileArgUpdate(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("POST /api/v2/runfile/args")

	if runfile.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	var req struct {
		Flag  string `json:"flag"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Runfile.Error(fmt.Sprintf("invalid request body: %v", err))
		writeJSONResponse(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if req.Flag == "" {
		logger.Runfile.Error("flag is required")
		writeJSONResponse(w, http.StatusBadRequest, nil, "flag is required")
		return
	}

	if err := runfile.SetArgValue(req.Flag, req.Value); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to set arg %s: %v", req.Flag, err))
		writeJSONResponse(w, http.StatusBadRequest, nil, fmt.Sprintf("failed to set arg: %v", err))
		return
	}

	logger.Runfile.Info(fmt.Sprintf("updated arg %s to %s", req.Flag, req.Value))
	writeJSONResponse(w, http.StatusOK, map[string]string{"flag": req.Flag, "value": req.Value}, "")
}

// HandleRunfile handles GET /api/v2/runfile
func HandleRunfile(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("GET /api/v2/runfile")

	if runfile.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	apiRunfile := toAPIRunFile(runfile.CurrentRunfile)
	logger.Runfile.Info("fetched runfile")
	writeJSONResponse(w, http.StatusOK, apiRunfile, "")
}

// HandleRunfileSave handles POST /api/v2/runfile/save
func HandleRunfileSave(w http.ResponseWriter, r *http.Request) {
	logger.Runfile.Debug("POST /api/v2/runfile/save")

	if runfile.CurrentRunfile == nil {
		logger.Runfile.Error("runfile not loaded")
		writeJSONResponse(w, http.StatusInternalServerError, nil, "runfile not loaded")
		return
	}

	if err := runfile.SaveRunfile(); err != nil {
		logger.Runfile.Error(fmt.Sprintf("failed to save runfile: %v", err))
		writeJSONResponse(w, http.StatusInternalServerError, nil, fmt.Sprintf("failed to save runfile: %v", err))
		return
	}

	logger.Runfile.Info("runfile saved")
	writeJSONResponse(w, http.StatusOK, "runfile saved", "")
}

// HandleSetRunfileGame reloads the runfile and restarts most of the server. It can also be used to reload the runfile from Disk as a hard reset.
func HandleSetRunfileGame(w http.ResponseWriter, r *http.Request) {

	reloadMu.Lock()
	defer reloadMu.Unlock()

	// Restrict to POST method
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read and validate request body
	var request struct {
		Game string `json:"game"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	game := strings.TrimSpace(request.Game)
	if game == "" {
		http.Error(w, "Game cannot be empty", http.StatusBadRequest)
		return
	}

	// Call InitRunfile to handle the runfile update
	if err := loader.InitRunfile(game); err != nil {
		logger.Core.Debug("Failed to initialize runfile: " + err.Error())
		http.Error(w, "Failed to initialize runfile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logger.Core.Info("Runfile game updated successfully to " + game)

	// Prepare response
	response := struct {
		Message string `json:"message"`
		Game    string `json:"game"`
	}{
		Message: "Monitor console for update status",
		Game:    game,
	}

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func HandleReloadRunfile(w http.ResponseWriter, r *http.Request) {
	logger.Web.Debug("Received reloadrunfile request from API")
	reloadMu.Lock()
	defer reloadMu.Unlock()
	// accept only GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	loader.ReloadRunfile()

	// Set response headers and write JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "OK"}); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
