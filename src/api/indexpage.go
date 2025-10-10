package api

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/config"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/localization"
	"github.com/JacksonTheMaster/StationeersServerUI/v5/src/logger"
)

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	htmlFS, err := fs.Sub(config.V1UIFS, "UIMod/onboard_bundled/ui")
	if err != nil {
		http.Error(w, "Error accessing Virt FS: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFS(htmlFS, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Core.Error("failed to serve v1 Index.html")
		return
	}

	var Identifier string

	if config.SSUIIdentifier == "" {
		Identifier = " (" + config.GetBranch() + ")"
	} else {
		Identifier = ": " + config.GetSSUIIdentifier()
	}

	data := IndexTemplateData{
		Version:                        config.GetVersion(),
		Branch:                         config.GetBranch(),
		SSUIIdentifier:                 Identifier,
		UIText_StartButton:             localization.GetString("UIText_StartButton"),
		UIText_StopButton:              localization.GetString("UIText_StopButton"),
		UIText_Settings:                localization.GetString("UIText_Settings"),
		UIText_Update_SteamCMD:         localization.GetString("UIText_Update_SteamCMD"),
		UIText_Console:                 localization.GetString("UIText_Console"),
		UIText_Detection_Events:        localization.GetString("UIText_Detection_Events"),
		UIText_Backend_Log:             localization.GetString("UIText_Backend_Log"),
		UIText_Backup_Manager:          localization.GetString("UIText_Backup_Manager"),
		UIText_Connected_PlayersHeader: localization.GetString("UIText_Connected_PlayersHeader"),
		UIText_Discord_Info:            localization.GetString("UIText_Discord_Info"),
		UIText_API_Info:                localization.GetString("UIText_API_Info"),
		UIText_Copyright1:              localization.GetString("UIText_Copyright1"),
		UIText_Copyright2:              localization.GetString("UIText_Copyright2"),
	}
	if data.Version == "" {
		data.Version = "unknown"
	}
	if data.Branch == "" {
		data.Branch = "unknown"
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
