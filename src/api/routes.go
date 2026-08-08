package api

import (
	"net/http"

	"github.com/SteamServerUI/SteamServerUI/v7/src/api/pluginsapi"
)

var GlobalWebProtectedMux *http.ServeMux

// SetupSocketAPIRoutes exposes the unsafe plugin registration surface only when
// the caller has already enabled the quarantined plugin runtime.
func SetupSocketAPIRoutes(apiMux *http.ServeMux) {
	apiMux.HandleFunc("/api/v3/plugins/log", pluginsapi.PluginLogHandler)
	apiMux.HandleFunc("/api/v3/plugins/register", func(w http.ResponseWriter, r *http.Request) {
		pluginsapi.RegisterPluginRouteHandler(w, r, apiMux, GlobalWebProtectedMux)
	})
}
