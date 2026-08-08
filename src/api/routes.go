package api

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/backupapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/httpauth"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/legacyapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pages"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginsapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/runfileapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/settingsapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sscmapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sseapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sysinfoapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/settings"
)

var GlobalWebProtectedMux *http.ServeMux

// SetupAPIRoutes sets up API routes used by B O T H the web and socket servers
func SetupAPIRoutes() (*http.ServeMux, *http.ServeMux) {

	// Set up handlers with auth middleware
	mux := http.NewServeMux() // Use a mux to apply middleware globally

	// Unprotected auth routes
	twoboxformAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/twoboxform")
	mux.Handle("/twoboxform/", http.StripPrefix("/twoboxform/", http.FileServer(http.FS(twoboxformAssetsFS))))
	mux.HandleFunc("/auth/login", httpauth.SessionLoginHandler)
	mux.HandleFunc("/auth/logout", httpauth.SessionLogoutHandler)
	mux.HandleFunc("/api/v2/auth/setup/status", httpauth.SetupStatusHandler)
	mux.HandleFunc("/api/v2/auth/setup/bootstrap", httpauth.BootstrapOwnerHandler)
	mux.HandleFunc("/login", pages.ServeTwoBoxFormTemplate)

	// Protected routes (wrapped with middleware)
	protectedMux := http.NewServeMux()
	GlobalWebProtectedMux = protectedMux

	// http file server for ./SSUI/config/files
	protectedMux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(filepath.Join(config.GetSSUIFolder(), "config", "files")))))

	legacyAssetsFS, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/assets")
	protectedMux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(legacyAssetsFS))))

	protectedMux.HandleFunc("/legacy/config", pages.ServeConfigPage)
	protectedMux.HandleFunc("/legacy/detectionmanager", pages.ServeDetectionManager)

	// Index page(s)
	protectedMux.HandleFunc("/legacy", pages.ServeIndex)
	protectedMux.HandleFunc("/", pages.ServeSvelteUI)

	// --- SVELTE UI ---
	svelteAssetsFS, _ := fs.Sub(config.V1UIFS, "SSUI/onboard_bundled/v2/assets")
	protectedMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(svelteAssetsFS))))
	protectedMux.HandleFunc("/api/v2/loader/reloadbackend", HandleReloadBackend)

	// SSE routes
	protectedMux.HandleFunc("/console", sseapi.GetLogOutput)
	protectedMux.HandleFunc("/events", sseapi.GetEventOutput)
	protectedMux.HandleFunc("/logs/debug", sseapi.GetDebugLogOutput)
	protectedMux.HandleFunc("/logs/info", sseapi.GetInfoLogOutput)
	protectedMux.HandleFunc("/logs/warn", sseapi.GetWarnLogOutput)
	protectedMux.HandleFunc("/logs/error", sseapi.GetErrorLogOutput)
	protectedMux.HandleFunc("/logs/backend", sseapi.GetBackendLogOutput)

	// Server Control
	protectedMux.HandleFunc("/start", legacyapi.StartServer)
	protectedMux.HandleFunc("/stop", legacyapi.StopServer)
	protectedMux.HandleFunc("/api/v2/server/start", legacyapi.StartServer) // TODO: should return json & get their own functions
	protectedMux.HandleFunc("/api/v2/server/stop", legacyapi.StopServer)   // TODO: should return json & get their own functions
	protectedMux.HandleFunc("/api/v2/server/status", GetGameServerRunState)
	protectedMux.HandleFunc("/api/v2/server/status/connectedplayers", legacyapi.HandleConnectedPlayersList)

	// Configuration
	protectedMux.HandleFunc("/api/v2/SSCM/run", sscmapi.HandleCommand)           // Command execution via SSCM (needs to be enable, config.IsSSCMEnabled)
	protectedMux.HandleFunc("/api/v2/SSCM/enabled", sscmapi.HandleIsSSCMEnabled) // Check if SSCM is enabled
	protectedMux.HandleFunc("/api/v2/steamcmd/run", HandleRunSteamCMD)           // Run SteamCMD

	// Custom Detections
	protectedMux.HandleFunc("/api/v2/custom-detections", detectionmgr.HandleCustomDetection)
	protectedMux.HandleFunc("/api/v2/custom-detections/delete/", detectionmgr.HandleDeleteCustomDetection)
	// Authentication
	protectedMux.HandleFunc("/changeuser", pages.ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/adduser", httpauth.RegisterUserHandler)        // user registration and change password
	protectedMux.HandleFunc("/api/v2/auth/setup/apikey", httpauth.RegisterAPIKeyHandler) // apikey registration and change password
	protectedMux.HandleFunc("/api/v2/auth/whoami", httpauth.WhoAmIHandler)
	protectedMux.HandleFunc("/api/v2/auth/session", httpauth.SessionInfoHandler)

	// Setup
	protectedMux.HandleFunc("/setup", pages.ServeTwoBoxFormTemplate)
	protectedMux.HandleFunc("/api/v2/auth/setup/register", httpauth.RegisterUserHandler) // user registration
	protectedMux.HandleFunc("/api/v2/auth/setup/finalize", httpauth.ActivateAuthHandler)

	// SteamServerUI

	// --- RUNFILE ---
	protectedMux.HandleFunc("/api/v2/runfile/groups", runfileapi.HandleRunfileGroups)
	protectedMux.HandleFunc("/api/v2/runfile/args", runfileapi.HandleRunfileArgs)
	protectedMux.HandleFunc("/api/v2/runfile/args/update", runfileapi.HandleRunfileArgUpdate)
	protectedMux.HandleFunc("/api/v2/runfile/args/getarg", runfileapi.HandleRunfileGetArg)
	protectedMux.HandleFunc("/api/v2/runfile/save", runfileapi.HandleRunfileSave)
	protectedMux.HandleFunc("/api/v2/runfile/hardreset", runfileapi.HandleSetRunfileGame)
	protectedMux.HandleFunc("/api/v2/runfile/meta", runfileapi.HandleRunfileGetMeta)
	// --- LOADER ---
	protectedMux.HandleFunc("/api/v2/loader/reloadrunfile", runfileapi.HandleReloadRunfile)
	// --- SETTINGS ---
	protectedMux.HandleFunc("/api/v2/settings", settings.HandleRetrieveSettings)
	protectedMux.HandleFunc("/api/v2/settings/save", settings.HandleSaveSetting)
	protectedMux.HandleFunc("/api/v2/settings/files/upload", settingsapi.HandleFileUpload)
	protectedMux.HandleFunc("/api/v2/settings/files/background/upload", settingsapi.HandleBackgroundUpload)
	protectedMux.HandleFunc("/api/v2/settings/files/tls/upload", settingsapi.HandleTLSCertUpload)
	// --- OS STATS ---
	protectedMux.HandleFunc("/api/v2/osstats", sysinfoapi.HandleGetOsStats)
	// --- RUNFILE GALLERY ---
	protectedMux.HandleFunc("/api/v2/gallery", runfileapi.GalleryHandler)
	protectedMux.HandleFunc("/api/v2/gallery/select", runfileapi.GallerySelectHandler)

	// --- PLUGIN GALLERY ---
	protectedMux.HandleFunc("/api/v2/plugingallery", pluginsapi.PluginGalleryHandler)
	protectedMux.HandleFunc("/api/v2/plugingallery/select", pluginsapi.PluginSelectHandler)

	// --- FILE MANAGEMENT ---
	protectedMux.HandleFunc("/api/v2/files", runfileapi.GetFileList)
	protectedMux.HandleFunc("/api/v2/files/get", runfileapi.GetFile)
	protectedMux.HandleFunc("/api/v2/files/save", runfileapi.SaveFile)

	// --- PLUGINS ---
	protectedMux.HandleFunc("/api/v2/plugins/list/apiroutes", pluginsapi.HandleListPluginAPIRoutes)
	protectedMux.HandleFunc("/api/v2/plugins/list/names", pluginsapi.HandleListPluginNames)
	protectedMux.HandleFunc("/api/v2/plugins/stop", pluginsapi.HandleStopPlugin)

	// --- BACKUP ---
	protectedMux.HandleFunc("/api/v2/backup/create", backupapi.HandleBackupCreate)
	protectedMux.HandleFunc("/api/v2/backup/list", backupapi.HandleBackupList)
	protectedMux.HandleFunc("/api/v2/backup/restore", backupapi.HandleBackupRestore)
	protectedMux.HandleFunc("/api/v2/backup/status", backupapi.HandleBackupStatus)

	return mux, protectedMux
}

// SetupSocketAPIRoutes adds routes that are E X C L U S I V E L Y available via sockets (if debug mode is enabled, these routes are added to the http api as well)
func SetupSocketAPIRoutes(APIMux *http.ServeMux) {
	APIMux.HandleFunc("/api/v2/plugins/log", pluginsapi.PluginLogHandler)
	APIMux.HandleFunc("/api/v2/plugins/register", func(w http.ResponseWriter, r *http.Request) {
		pluginsapi.RegisterPluginRouteHandler(w, r, APIMux, GlobalWebProtectedMux)
	})
}
