package httpauth

import (
	"net/http"
	"sort"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/middleware"
	"github.com/SteamServerUI/SteamServerUI/v7/src/config"
)

func AuditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.WriteJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Method not allowed")
		return
	}
	events := append([]config.IdentityAuditEvent(nil), config.GetIdentityConfig().Audit...)
	sort.Slice(events, func(i, j int) bool { return events[i].CreatedAt.After(events[j].CreatedAt) })
	writeJSON(w, http.StatusOK, map[string]any{"events": events})
}
