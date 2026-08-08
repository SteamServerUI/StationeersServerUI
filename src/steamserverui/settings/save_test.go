package settings

import "testing"

func TestSettingEffects(t *testing.T) {
	if got := settingEffects("BackendEndpointPort"); len(got) != 1 || got[0] != "backend.reload" {
		t.Fatalf("BackendEndpointPort effects = %#v", got)
	}
	if got := settingEffects("IsDiscordEnabled"); len(got) != 1 || got[0] != "immediate" {
		t.Fatalf("IsDiscordEnabled effects = %#v", got)
	}
}
