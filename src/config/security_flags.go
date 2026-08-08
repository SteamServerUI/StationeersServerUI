package config

import (
	"os"
	"strings"
)

func GetPluginsEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("SSUI_ENABLE_UNSAFE_PLUGINS")), "true")
}
