package httpauth

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

type identityCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SetupStatusHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.ApplyCORS(w, r) {
		return
	}
	if r.Method != http.MethodGet {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"setupRequired": config.GetIdentityConfig().SetupRequired})
}

func BootstrapOwnerHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.ApplyCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	var request struct {
		SetupSecret string `json:"setupSecret"`
		Username    string `json:"username"`
		Password    string `json:"password"`
	}
	if err := decodeJSON(r, &request); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	owner, err := security.BootstrapOwner(request.SetupSecret, request.Username, request.Password, time.Now())
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "setup_failed", err.Error())
		return
	}
	respondWithNewSession(w, owner)
}

func SessionLoginHandler(w http.ResponseWriter, r *http.Request) {
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
	respondWithNewSession(w, user)
}

func SessionLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if middleware.ApplyCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	cookie, err := r.Cookie(security.SessionCookieName)
	if err == nil {
		_, session, authErr := security.AuthenticateSession(cookie.Value, time.Now())
		if authErr == nil && security.ValidateSessionCSRF(session, r.Header.Get("X-SSUI-CSRF")) {
			_ = security.RevokeSession(session.ID)
		}
	}
	clearSessionCookie(w)
	w.WriteHeader(http.StatusNoContent)
}

func SessionInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	principal, ok := middleware.PrincipalFromContext(r.Context())
	if !ok {
		middleware.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
		return
	}
	identity := config.GetIdentityConfig()
	user, ok := identity.Users[principal.UserID]
	if !ok {
		middleware.WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "User no longer exists")
		return
	}
	response := sessionResponse(user, principal)
	if principal.Credential == "session" {
		if cookie, err := r.Cookie(security.SessionCookieName); err == nil {
			response["csrf"] = security.CSRFForSessionCredential(cookie.Value)
		}
		if session, ok := middleware.SessionFromContext(r.Context()); ok {
			response["expiresAt"] = session.AbsoluteExpiresAt
		}
	}
	writeJSON(w, http.StatusOK, response)
}

func respondWithNewSession(w http.ResponseWriter, user config.IdentityUser) {
	session, credential, err := security.CreateSession(user.ID, time.Now())
	if err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "session_failed", "Could not create session")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     security.SessionCookieName,
		Value:    credential.Value,
		Path:     "/",
		Expires:  session.AbsoluteExpiresAt,
		MaxAge:   int(time.Until(session.AbsoluteExpiresAt).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	identity := config.GetIdentityConfig()
	principal := security.Principal{
		UserID:      user.ID,
		Username:    user.Username,
		Credential:  "session",
		Permissions: security.ResolvePermissions(user, identity),
	}
	response := sessionResponse(user, principal)
	response["csrf"] = credential.CSRF
	response["expiresAt"] = session.AbsoluteExpiresAt
	writeJSON(w, http.StatusOK, response)
}

func sessionResponse(user config.IdentityUser, principal security.Principal) map[string]any {
	permissions := make([]string, 0, len(principal.Permissions))
	for _, permission := range security.AllPermissions {
		if principal.Permissions[permission] {
			permissions = append(permissions, permission)
		}
	}
	return map[string]any{
		"user": map[string]any{
			"id":       user.ID,
			"username": user.Username,
			"groupIds": user.GroupIDs,
		},
		"credentialType": principal.Credential,
		"permissions":    permissions,
	}
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     security.SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(1, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func decodeJSON(r *http.Request, target any) error {
	if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return &requestError{"Content-Type must be application/json"}
	}
	decoder := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 64*1024))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return &requestError{"Invalid JSON body"}
	}
	return nil
}

type requestError struct{ message string }

func (err *requestError) Error() string { return err.message }

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
