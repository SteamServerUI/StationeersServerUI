package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
	"github.com/SteamServerUI/SteamServerUI/v7/src/core/security"
)

type identityContextKey string

const (
	principalContextKey identityContextKey = "principal"
	sessionContextKey   identityContextKey = "session"
)

func IdentityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.GetIdentityConfig().SetupRequired {
			WriteJSONError(w, http.StatusServiceUnavailable, "setup_required", "Owner setup is required")
			return
		}

		principal, session, err := authenticateRequest(r)
		if err != nil {
			WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
			return
		}
		if principal.Credential == "session" && changesState(r.Method) {
			if !sameOrigin(r) || !security.ValidateSessionCSRF(session, r.Header.Get("X-SSUI-CSRF")) {
				WriteJSONError(w, http.StatusForbidden, "csrf_failed", "CSRF validation failed")
				return
			}
		}

		ctx := context.WithValue(r.Context(), principalContextKey, principal)
		if principal.Credential == "session" {
			ctx = context.WithValue(ctx, sessionContextKey, session)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequirePermission(permission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			WriteJSONError(w, http.StatusUnauthorized, "unauthorized", "Authentication required")
			return
		}
		if !principal.Permissions[permission] {
			WriteJSONError(w, http.StatusForbidden, "forbidden", "Permission denied")
			return
		}
		next(w, r)
	}
}

func PrincipalFromContext(ctx context.Context) (security.Principal, bool) {
	principal, ok := ctx.Value(principalContextKey).(security.Principal)
	return principal, ok
}

func SessionFromContext(ctx context.Context) (config.IdentitySession, bool) {
	session, ok := ctx.Value(sessionContextKey).(config.IdentitySession)
	return session, ok
}

func WriteJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func ApplyCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin != "" && sameOrigin(r) {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Vary", "Origin")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-SSUI-CSRF")
	if r.Method == http.MethodOptions {
		if origin != "" && !sameOrigin(r) {
			WriteJSONError(w, http.StatusForbidden, "origin_denied", "Origin is not allowed")
		} else {
			w.WriteHeader(http.StatusNoContent)
		}
		return true
	}
	return false
}

func authenticateRequest(r *http.Request) (security.Principal, config.IdentitySession, error) {
	authorization := r.Header.Get("Authorization")
	if strings.HasPrefix(authorization, "Bearer ") {
		principal, _, err := security.AuthenticateToken(strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer ")), time.Now())
		return principal, config.IdentitySession{}, err
	}
	cookie, err := r.Cookie(security.SessionCookieName)
	if err != nil {
		return security.Principal{}, config.IdentitySession{}, err
	}
	return security.AuthenticateSession(cookie.Value, time.Now())
}

func sameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return origin == scheme+"://"+r.Host
}

func changesState(method string) bool {
	return method != http.MethodGet && method != http.MethodHead && method != http.MethodOptions
}
