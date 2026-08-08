package httpauth

import (
	"net/http"
	"sort"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

func DesktopLoginHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.ApplyCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	var request identityCredentials
	if err := decodeJSON(r, &request); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	user, err := security.AuthenticateUser(request.Username, request.Password)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusUnauthorized, "invalid_credentials", "Invalid username or password")
		return
	}
	identity := config.GetIdentityConfig()
	grants := security.ResolvePermissions(user, identity)
	scopes := make([]string, 0, len(grants))
	for permission := range grants {
		scopes = append(scopes, permission)
	}
	sort.Strings(scopes)
	expiresAt := time.Now().AddDate(0, 6, 0)
	_, secret, err := security.CreateNamedToken(user.ID, "SSUI Desktop", scopes, &expiresAt, time.Now())
	if err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "token_failed", "Could not create desktop credential")
		return
	}
	principal := security.Principal{
		UserID:      user.ID,
		Username:    user.Username,
		Credential:  "token",
		Permissions: grants,
	}
	response := sessionResponse(user, principal)
	response["token"] = secret
	response["expiresAt"] = expiresAt
	writeJSON(w, http.StatusOK, response)
}

func DesktopLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok || principal.Credential != "token" {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_credential", "A desktop token is required")
		return
	}
	if err := security.RevokeTokenAs(principal.CredentialID, principal.UserID, principal.Username, time.Now()); err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "revoke_failed", "Could not revoke desktop token")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
