package httpauth

import (
	"errors"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
	"github.com/google/uuid"
)

type userView struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Enabled   bool      `json:"enabled"`
	GroupIDs  []string  `json:"groupIds"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsers(w)
	case http.MethodPost:
		createUser(w, r)
	default:
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v2/auth/users/")
	if id == "" || strings.Contains(id, "/") {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}
	var request struct {
		Username *string   `json:"username"`
		Password *string   `json:"password"`
		Enabled  *bool     `json:"enabled"`
		GroupIDs *[]string `json:"groupIds"`
	}
	if err := decodeJSON(r, &request); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	principal, _ := middleware.PrincipalFromContext(r.Context())
	identity := config.GetIdentityConfig()
	current, ok := identity.Users[id]
	if !ok {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	updated := current
	if request.Username != nil {
		if err := security.ValidateUsername(*request.Username); err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_username", err.Error())
			return
		}
		updated.Username = strings.TrimSpace(*request.Username)
		updated.Normalized = security.NormalizeUsername(*request.Username)
		if usernameExists(identity, updated.Normalized, id) {
			middleware.WriteJSONError(w, http.StatusConflict, "username_exists", "Username already exists")
			return
		}
	}
	if request.Enabled != nil {
		updated.Enabled = *request.Enabled
	}
	if request.GroupIDs != nil {
		groups, err := validateGroupIDs(identity, *request.GroupIDs)
		if err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_groups", err.Error())
			return
		}
		if err := canAssignGroups(identity, groups, principal); err != nil {
			middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", err.Error())
			return
		}
		if ownerMembershipChanged(current.GroupIDs, groups) && !principal.Permissions[security.PermissionSecurityManage] {
			middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", "Changing ownership requires security.manage")
			return
		}
		updated.GroupIDs = groups
	}
	if isEnabledOwner(current) && !isEnabledOwner(updated) && enabledOwnerCount(identity) == 1 {
		middleware.WriteJSONError(w, http.StatusConflict, "last_owner", "The last enabled owner cannot be removed or disabled")
		return
	}
	if request.Password != nil {
		hash, err := security.HashIdentityPassword(*request.Password)
		if err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_password", err.Error())
			return
		}
		updated.PasswordHash = hash
	}
	updated.UpdatedAt = time.Now()
	if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
		value.Users[id] = updated
		if request.Password != nil || !updated.Enabled {
			removeUserCredentials(value, id)
		}
		return nil
	}); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not save user")
		return
	}
	writeJSON(w, http.StatusOK, toUserView(updated))
}

func GroupsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		identity := config.GetIdentityConfig()
		groups := make([]config.IdentityGroup, 0, len(identity.Groups))
		for _, group := range identity.Groups {
			groups = append(groups, group)
		}
		sort.Slice(groups, func(i, j int) bool { return strings.ToLower(groups[i].Name) < strings.ToLower(groups[j].Name) })
		writeJSON(w, http.StatusOK, map[string]any{"groups": groups, "permissions": security.AllPermissions})
	case http.MethodPost:
		createGroup(w, r)
	default:
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func GroupHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v2/auth/groups/")
	if id == "" || strings.Contains(id, "/") {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "Group not found")
		return
	}
	identity := config.GetIdentityConfig()
	group, ok := identity.Groups[id]
	if !ok {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "Group not found")
		return
	}
	if group.System {
		middleware.WriteJSONError(w, http.StatusConflict, "system_group", "System groups cannot be changed")
		return
	}
	switch r.Method {
	case http.MethodPut:
		var request struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Permissions []string `json:"permissions"`
		}
		if err := decodeJSON(r, &request); err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		principal, _ := middleware.PrincipalFromContext(r.Context())
		permissions, err := validatePermissionGrant(principal, request.Permissions)
		if err != nil {
			middleware.WriteJSONError(w, http.StatusForbidden, "invalid_permissions", err.Error())
			return
		}
		name := strings.TrimSpace(request.Name)
		if name == "" || groupNameExists(identity, name, id) {
			middleware.WriteJSONError(w, http.StatusConflict, "group_name_exists", "Group name is empty or already exists")
			return
		}
		group.Name = name
		group.Description = strings.TrimSpace(request.Description)
		group.Permissions = permissions
		group.UpdatedAt = time.Now()
		if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
			value.Groups[id] = group
			return nil
		}); err != nil {
			middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not save group")
			return
		}
		writeJSON(w, http.StatusOK, group)
	case http.MethodDelete:
		if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
			delete(value.Groups, id)
			for userID, user := range value.Users {
				user.GroupIDs = withoutValue(user.GroupIDs, id)
				value.Users[userID] = user
			}
			return nil
		}); err != nil {
			middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not delete group")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func TokensHandler(w http.ResponseWriter, r *http.Request) {
	principal, _ := middleware.PrincipalFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		identity := config.GetIdentityConfig()
		tokens := make([]config.IdentityToken, 0)
		for _, token := range identity.Tokens {
			if token.OwnerID == principal.UserID || principal.Permissions[security.PermissionAllTokensManage] {
				token.SecretHash = ""
				tokens = append(tokens, token)
			}
		}
		writeJSON(w, http.StatusOK, map[string]any{"tokens": tokens})
	case http.MethodPost:
		var request struct {
			Name      string     `json:"name"`
			Scopes    []string   `json:"scopes"`
			ExpiresAt *time.Time `json:"expiresAt"`
		}
		if err := decodeJSON(r, &request); err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
		token, secret, err := security.CreateNamedToken(principal.UserID, request.Name, request.Scopes, request.ExpiresAt, time.Now())
		if err != nil {
			middleware.WriteJSONError(w, http.StatusBadRequest, "token_failed", err.Error())
			return
		}
		token.SecretHash = ""
		writeJSON(w, http.StatusCreated, map[string]any{"token": token, "secret": secret})
	default:
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
	}
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v2/auth/tokens/")
	principal, _ := middleware.PrincipalFromContext(r.Context())
	token, ok := config.GetIdentityConfig().Tokens[id]
	if !ok {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "Token not found")
		return
	}
	if token.OwnerID != principal.UserID && !principal.Permissions[security.PermissionAllTokensManage] {
		middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}
	if err := security.RevokeToken(id, time.Now()); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not revoke token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func SessionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	principal, _ := middleware.PrincipalFromContext(r.Context())
	identity := config.GetIdentityConfig()
	sessions := make([]config.IdentitySession, 0)
	for _, session := range identity.Sessions {
		if session.UserID == principal.UserID || principal.Permissions[security.PermissionSessionsManage] {
			session.SecretHash = ""
			session.CSRFHash = ""
			sessions = append(sessions, session)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/api/v2/auth/sessions/")
	principal, _ := middleware.PrincipalFromContext(r.Context())
	session, ok := config.GetIdentityConfig().Sessions[id]
	if !ok {
		middleware.WriteJSONError(w, http.StatusNotFound, "not_found", "Session not found")
		return
	}
	if session.UserID != principal.UserID && !principal.Permissions[security.PermissionSessionsManage] {
		middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", "Permission denied")
		return
	}
	if err := security.RevokeSession(id); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not revoke session")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func listUsers(w http.ResponseWriter) {
	identity := config.GetIdentityConfig()
	users := make([]userView, 0, len(identity.Users))
	for _, user := range identity.Users {
		users = append(users, toUserView(user))
	}
	sort.Slice(users, func(i, j int) bool { return users[i].Username < users[j].Username })
	writeJSON(w, http.StatusOK, map[string]any{"users": users})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		GroupIDs []string `json:"groupIds"`
	}
	if err := decodeJSON(r, &request); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	if err := security.ValidateUsername(request.Username); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_username", err.Error())
		return
	}
	hash, err := security.HashIdentityPassword(request.Password)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_password", err.Error())
		return
	}
	identity := config.GetIdentityConfig()
	if usernameExists(identity, security.NormalizeUsername(request.Username), "") {
		middleware.WriteJSONError(w, http.StatusConflict, "username_exists", "Username already exists")
		return
	}
	groups, err := validateGroupIDs(identity, request.GroupIDs)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_groups", err.Error())
		return
	}
	principal, _ := middleware.PrincipalFromContext(r.Context())
	if err := canAssignGroups(identity, groups, principal); err != nil {
		middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", err.Error())
		return
	}
	if contains(groups, security.OwnerGroupID) && !principal.Permissions[security.PermissionSecurityManage] {
		middleware.WriteJSONError(w, http.StatusForbidden, "forbidden", "Creating an owner requires security.manage")
		return
	}
	now := time.Now()
	user := config.IdentityUser{
		ID:           uuid.NewString(),
		Username:     strings.TrimSpace(request.Username),
		Normalized:   security.NormalizeUsername(request.Username),
		PasswordHash: hash,
		Enabled:      true,
		GroupIDs:     groups,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
		value.Users[user.ID] = user
		return nil
	}); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not save user")
		return
	}
	writeJSON(w, http.StatusCreated, toUserView(user))
}

func createGroup(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
	}
	if err := decodeJSON(r, &request); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	principal, _ := middleware.PrincipalFromContext(r.Context())
	permissions, err := validatePermissionGrant(principal, request.Permissions)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusForbidden, "invalid_permissions", err.Error())
		return
	}
	name := strings.TrimSpace(request.Name)
	identity := config.GetIdentityConfig()
	if name == "" || groupNameExists(identity, name, "") {
		middleware.WriteJSONError(w, http.StatusConflict, "group_name_exists", "Group name is empty or already exists")
		return
	}
	now := time.Now()
	group := config.IdentityGroup{
		ID:          uuid.NewString(),
		Name:        name,
		Description: strings.TrimSpace(request.Description),
		Permissions: permissions,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := config.MutateIdentityConfig(func(value *config.IdentityConfig) error {
		value.Groups[group.ID] = group
		return nil
	}); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "save_failed", "Could not save group")
		return
	}
	writeJSON(w, http.StatusCreated, group)
}

func validatePermissionGrant(principal security.Principal, requested []string) ([]string, error) {
	permissions := make([]string, 0, len(requested))
	seen := make(map[string]bool)
	for _, permission := range requested {
		if !security.IsPermission(permission) {
			return nil, errors.New("unknown permission: " + permission)
		}
		if !principal.Permissions[permission] {
			return nil, errors.New("permission exceeds your access: " + permission)
		}
		if !seen[permission] {
			seen[permission] = true
			permissions = append(permissions, permission)
		}
	}
	return permissions, nil
}

func validateGroupIDs(identity config.IdentityConfig, requested []string) ([]string, error) {
	groups := make([]string, 0, len(requested))
	seen := make(map[string]bool)
	for _, id := range requested {
		if _, ok := identity.Groups[id]; !ok {
			return nil, errors.New("group does not exist: " + id)
		}
		if !seen[id] {
			seen[id] = true
			groups = append(groups, id)
		}
	}
	return groups, nil
}

func canAssignGroups(identity config.IdentityConfig, groups []string, principal security.Principal) error {
	for _, groupID := range groups {
		group := identity.Groups[groupID]
		for _, permission := range group.Permissions {
			if !principal.Permissions[permission] {
				return errors.New("group exceeds your access: " + group.Name)
			}
		}
	}
	return nil
}

func usernameExists(identity config.IdentityConfig, normalized, exceptID string) bool {
	for id, user := range identity.Users {
		if id != exceptID && user.Normalized == normalized {
			return true
		}
	}
	return false
}

func groupNameExists(identity config.IdentityConfig, name, exceptID string) bool {
	for id, group := range identity.Groups {
		if id != exceptID && strings.EqualFold(group.Name, name) {
			return true
		}
	}
	return false
}

func enabledOwnerCount(identity config.IdentityConfig) int {
	count := 0
	for _, user := range identity.Users {
		if isEnabledOwner(user) {
			count++
		}
	}
	return count
}

func isEnabledOwner(user config.IdentityUser) bool {
	return user.Enabled && contains(user.GroupIDs, security.OwnerGroupID)
}

func ownerMembershipChanged(before, after []string) bool {
	return contains(before, security.OwnerGroupID) != contains(after, security.OwnerGroupID)
}

func removeUserCredentials(identity *config.IdentityConfig, userID string) {
	for id, session := range identity.Sessions {
		if session.UserID == userID {
			delete(identity.Sessions, id)
		}
	}
	for id, token := range identity.Tokens {
		if token.OwnerID == userID {
			delete(identity.Tokens, id)
		}
	}
}

func toUserView(user config.IdentityUser) userView {
	return userView{
		ID:        user.ID,
		Username:  user.Username,
		Enabled:   user.Enabled,
		GroupIDs:  user.GroupIDs,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func withoutValue(values []string, unwanted string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value != unwanted {
			result = append(result, value)
		}
	}
	return result
}
