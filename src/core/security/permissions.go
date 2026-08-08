package security

const (
	PermissionServerView      = "server.view"
	PermissionServerControl   = "server.control"
	PermissionConsoleWrite    = "console.write"
	PermissionLogsView        = "logs.view"
	PermissionBackupsView     = "backups.view"
	PermissionBackupsCreate   = "backups.create"
	PermissionBackupsRestore  = "backups.restore"
	PermissionBackupsDelete   = "backups.delete"
	PermissionRunfilesView    = "runfiles.view"
	PermissionRunfilesManage  = "runfiles.manage"
	PermissionSettingsView    = "settings.view"
	PermissionSettingsManage  = "settings.manage"
	PermissionFilesRead       = "files.read"
	PermissionFilesWrite      = "files.write"
	PermissionSteamCMDRun     = "steamcmd.run"
	PermissionSSCMRun         = "sscm.run"
	PermissionPluginsView     = "plugins.view"
	PermissionPluginsManage   = "plugins.manage"
	PermissionUsersView       = "users.view"
	PermissionUsersManage     = "users.manage"
	PermissionGroupsView      = "groups.view"
	PermissionGroupsManage    = "groups.manage"
	PermissionOwnTokensManage = "tokens.own.manage"
	PermissionAllTokensManage = "tokens.all.manage"
	PermissionSessionsManage  = "sessions.manage"
	PermissionAuditView       = "audit.view"
	PermissionBackendReload   = "backend.reload"
	PermissionBackendUpdate   = "backend.update"
	PermissionSecurityManage  = "security.manage"
)

var AllPermissions = []string{
	PermissionServerView,
	PermissionServerControl,
	PermissionConsoleWrite,
	PermissionLogsView,
	PermissionBackupsView,
	PermissionBackupsCreate,
	PermissionBackupsRestore,
	PermissionBackupsDelete,
	PermissionRunfilesView,
	PermissionRunfilesManage,
	PermissionSettingsView,
	PermissionSettingsManage,
	PermissionFilesRead,
	PermissionFilesWrite,
	PermissionSteamCMDRun,
	PermissionSSCMRun,
	PermissionPluginsView,
	PermissionPluginsManage,
	PermissionUsersView,
	PermissionUsersManage,
	PermissionGroupsView,
	PermissionGroupsManage,
	PermissionOwnTokensManage,
	PermissionAllTokensManage,
	PermissionSessionsManage,
	PermissionAuditView,
	PermissionBackendReload,
	PermissionBackendUpdate,
	PermissionSecurityManage,
}

func IsPermission(value string) bool {
	for _, permission := range AllPermissions {
		if permission == value {
			return true
		}
	}
	return false
}
