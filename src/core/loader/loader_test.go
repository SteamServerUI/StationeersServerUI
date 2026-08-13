package loader

import "testing"

func TestReloadBackendReloadsDiscordBot(t *testing.T) {
	original := reloadDiscordBotFunc
	defer func() {
		reloadDiscordBotFunc = original
	}()

	called := false
	reloadDiscordBotFunc = func() {
		called = true
	}

	ReloadBackend()

	if !called {
		t.Fatal("expected ReloadBackend to reload Discord bot")
	}
}
