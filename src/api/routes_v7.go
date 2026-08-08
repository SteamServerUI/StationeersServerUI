package api

import (
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/backupapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/httpauth"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/legacyapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pages"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginsapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/runfileapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/settingsapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sscmapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sseapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/api/sysinfoapi"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/SteamServerUI/SteamServerUI/v7/src/managers/detectionmgr"
	"github.com/SteamServerUI/SteamServerUI/v7/src/steamserverui/settings"
)

func SetupV7APIRoutes() (*http.ServeMux, *http.ServeMux) {
	public := http.NewServeMux()
	twoBoxAssets, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/twoboxform")
	public.Handle("/twoboxform/", http.StripPrefix("/twoboxform/", http.FileServer(http.FS(twoBoxAssets))))
	public.HandleFunc("/auth/login", httpauth.SessionLoginHandler)
	public.HandleFunc("/auth/logout", httpauth.SessionLogoutHandler)
	public.HandleFunc("/api/v2/auth/setup/status", httpauth.SetupStatusHandler)
	public.HandleFunc("/api/v2/auth/setup/bootstrap", httpauth.BootstrapOwnerHandler)

	protected := http.NewServeMux()
	GlobalWebProtectedMux = protected

	protected.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(filepath.Join(config.GetSSUIFolder(), "config", "files")))))
	legacyAssets, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/assets")
	protected.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(legacyAssets))))
	svelteAssets, _ := fs.Sub(config.GetV1UIFS(), "SSUI/onboard_bundled/v2/assets")
	protected.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(svelteAssets))))
	protected.HandleFunc("/legacy/config", pages.ServeConfigPage)
	protected.HandleFunc("/legacy/detectionmanager", pages.ServeDetectionManager)
	protected.HandleFunc("/legacy", pages.ServeIndex)
	protected.HandleFunc("/", pages.ServeSvelteUI)

	protect(protected, "/api/v2/loader/reloadbackend", security.PermissionBackendReload, HandleReloadBackend)
	protect(protected, "/console", security.PermissionLogsView, sseapi.GetLogOutput)
	protect(protected, "/events", security.PermissionLogsView, sseapi.GetEventOutput)
	protect(protected, "/logs/debug", security.PermissionLogsView, sseapi.GetDebugLogOutput)
	protect(protected, "/logs/info", security.PermissionLogsView, sseapi.GetInfoLogOutput)
	protect(protected, "/logs/warn", security.PermissionLogsView, sseapi.GetWarnLogOutput)
	protect(protected, "/logs/error", security.PermissionLogsView, sseapi.GetErrorLogOutput)
	protect(protected, "/logs/backend", security.PermissionLogsView, sseapi.GetBackendLogOutput)

	protect(protected, "/start", security.PermissionServerControl, legacyapi.StartServer)
	protect(protected, "/stop", security.PermissionServerControl, legacyapi.StopServer)
	protect(protected, "/api/v2/server/start", security.PermissionServerControl, legacyapi.StartServer)
	protect(protected, "/api/v2/server/stop", security.PermissionServerControl, legacyapi.StopServer)
	protect(protected, "/api/v2/server/status", security.PermissionServerView, GetGameServerRunState)
	protect(protected, "/api/v2/server/status/connectedplayers", security.PermissionServerView, legacyapi.HandleConnectedPlayersList)

	protect(protected, "/api/v2/SSCM/run", security.PermissionSSCMRun, sscmapi.HandleCommand)
	protect(protected, "/api/v2/SSCM/enabled", security.PermissionSettingsView, sscmapi.HandleIsSSCMEnabled)
	protect(protected, "/api/v2/steamcmd/run", security.PermissionSteamCMDRun, HandleRunSteamCMD)
	protect(protected, "/api/v2/custom-detections", security.PermissionSettingsManage, detectionmgr.HandleCustomDetection)
	protect(protected, "/api/v2/custom-detections/delete/", security.PermissionSettingsManage, detectionmgr.HandleDeleteCustomDetection)

	protected.HandleFunc("/api/v2/auth/session", httpauth.SessionInfoHandler)
	protected.HandleFunc("/api/v2/auth/whoami", httpauth.SessionInfoHandler)
	protect(protected, "/api/v2/auth/users", security.PermissionUsersManage, httpauth.UsersHandler)
	protect(protected, "/api/v2/auth/users/", security.PermissionUsersManage, httpauth.UserHandler)
	protect(protected, "/api/v2/auth/groups", security.PermissionGroupsManage, httpauth.GroupsHandler)
	protect(protected, "/api/v2/auth/groups/", security.PermissionGroupsManage, httpauth.GroupHandler)
	protect(protected, "/api/v2/auth/tokens", security.PermissionOwnTokensManage, httpauth.TokensHandler)
	protect(protected, "/api/v2/auth/tokens/", security.PermissionOwnTokensManage, httpauth.TokenHandler)
	protected.HandleFunc("/api/v2/auth/sessions", httpauth.SessionsHandler)
	protected.HandleFunc("/api/v2/auth/sessions/", httpauth.SessionHandler)

	protect(protected, "/api/v2/runfile/groups", security.PermissionRunfilesView, runfileapi.HandleRunfileGroups)
	protect(protected, "/api/v2/runfile/args", security.PermissionRunfilesView, runfileapi.HandleRunfileArgs)
	protect(protected, "/api/v2/runfile/args/getarg", security.PermissionRunfilesView, runfileapi.HandleRunfileGetArg)
	protect(protected, "/api/v2/runfile/meta", security.PermissionRunfilesView, runfileapi.HandleRunfileGetMeta)
	protect(protected, "/api/v2/runfile/args/update", security.PermissionRunfilesManage, runfileapi.HandleRunfileArgUpdate)
	protect(protected, "/api/v2/runfile/save", security.PermissionRunfilesManage, runfileapi.HandleRunfileSave)
	protect(protected, "/api/v2/runfile/hardreset", security.PermissionRunfilesManage, runfileapi.HandleSetRunfileGame)
	protect(protected, "/api/v2/loader/reloadrunfile", security.PermissionRunfilesManage, runfileapi.HandleReloadRunfile)

	protect(protected, "/api/v2/settings", security.PermissionSettingsView, settings.HandleRetrieveSettings)
	protect(protected, "/api/v2/settings/save", security.PermissionSettingsManage, settings.HandleSaveSetting)
	protect(protected, "/api/v2/settings/files/upload", security.PermissionFilesWrite, settingsapi.HandleFileUpload)
	protect(protected, "/api/v2/settings/files/background/upload", security.PermissionSettingsManage, settingsapi.HandleBackgroundUpload)
	protect(protected, "/api/v2/settings/files/tls/upload", security.PermissionSecurityManage, settingsapi.HandleTLSCertUpload)
	protect(protected, "/api/v2/osstats", security.PermissionServerView, sysinfoapi.HandleGetOsStats)

	protect(protected, "/api/v2/gallery", security.PermissionRunfilesView, runfileapi.GalleryHandler)
	protect(protected, "/api/v2/gallery/select", security.PermissionRunfilesManage, runfileapi.GallerySelectHandler)
	protect(protected, "/api/v2/plugingallery", security.PermissionPluginsView, pluginsapi.PluginGalleryHandler)
	protect(protected, "/api/v2/plugingallery/select", security.PermissionPluginsManage, pluginsapi.PluginSelectHandler)
	protect(protected, "/api/v2/files", security.PermissionFilesRead, runfileapi.GetFileList)
	protect(protected, "/api/v2/files/get", security.PermissionFilesRead, runfileapi.GetFile)
	protect(protected, "/api/v2/files/save", security.PermissionFilesWrite, runfileapi.SaveFile)
	protect(protected, "/api/v2/plugins/list/apiroutes", security.PermissionPluginsView, pluginsapi.HandleListPluginAPIRoutes)
	protect(protected, "/api/v2/plugins/list/names", security.PermissionPluginsView, pluginsapi.HandleListPluginNames)
	protect(protected, "/api/v2/plugins/stop", security.PermissionPluginsManage, pluginsapi.HandleStopPlugin)

	protect(protected, "/api/v2/backup/create", security.PermissionBackupsCreate, backupapi.HandleBackupCreate)
	protect(protected, "/api/v2/backup/list", security.PermissionBackupsView, backupapi.HandleBackupList)
	protect(protected, "/api/v2/backup/restore", security.PermissionBackupsRestore, backupapi.HandleBackupRestore)
	protect(protected, "/api/v2/backup/status", security.PermissionBackupsView, backupapi.HandleBackupStatus)
	return public, protected
}

func protect(mux *http.ServeMux, path, permission string, handler http.HandlerFunc) {
	mux.HandleFunc(path, middleware.RequirePermission(permission, handler))
}
