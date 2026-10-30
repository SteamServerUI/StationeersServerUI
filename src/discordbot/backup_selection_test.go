package discordbot

import (
	"strings"
	"testing"
	"time"

	"github.com/SteamServerUI/StationeersServerUI/v6/src/config"
	"github.com/bwmarrin/discordgo"
)

func TestBackupCommandRequiresName(t *testing.T) {
	name := "nested/040926_112440_auto.save"
	option := &discordgo.ApplicationCommandInteractionDataOption{Name: "name", Type: discordgo.ApplicationCommandOptionString, Value: name}
	if got, err := backupCommandName([]*discordgo.ApplicationCommandInteractionDataOption{option}); err != nil || got != name {
		t.Fatalf("filename option: %q, %v", got, err)
	}
	for _, options := range [][]*discordgo.ApplicationCommandInteractionDataOption{
		nil, {nil}, {option, option},
		{{Name: "index", Type: discordgo.ApplicationCommandOptionInteger, Value: float64(17)}},
		{{Name: "index", Type: discordgo.ApplicationCommandOptionString, Value: "17"}},
		{{Name: "name", Type: discordgo.ApplicationCommandOptionString, Value: float64(17)}},
		{{Name: "name", Type: discordgo.ApplicationCommandOptionString, Value: ""}},
	} {
		if got, err := backupCommandName(options); err == nil {
			t.Fatalf("accepted stale or malformed options as %q", got)
		}
	}
}

func TestActiveRestoreMenuRetainsItsBackupName(t *testing.T) {
	config.ConfigMu.Lock()
	enabled := config.DiscordRestoreVoteEnabled
	config.DiscordRestoreVoteEnabled = true
	config.ConfigMu.Unlock()
	t.Cleanup(func() {
		config.ConfigMu.Lock()
		config.DiscordRestoreVoteEnabled = enabled
		config.ConfigMu.Unlock()
		resetDiscordVotes()
	})
	resetDiscordVotes()
	discordVotes.Lock()
	discordVotes.restore = &activeVote{target: restoreVoteTarget{Name: "first.save"}, required: 3, voters: map[string]struct{}{"a": {}}, expiresAt: time.Now().Add(time.Minute)}
	discordVotes.Unlock()
	menu := buildVoteMenuOptions()
	var selection string
	for _, option := range menu {
		if strings.HasPrefix(option.Value, "restore:") {
			selection = option.Value
		}
	}
	if selection != "restore:first.save" {
		t.Fatalf("active menu did not pin its name: %q", selection)
	}
	// A later vote must not reinterpret a menu that was shown for the first save.
	discordVotes.Lock()
	discordVotes.restore.target.Name = "second.save"
	discordVotes.Unlock()
	result := castDiscordVote(voteRestore, restoreVoteTarget{Name: strings.TrimPrefix(selection, "restore:")}, "b")
	if !strings.Contains(result.message, "different backup") {
		t.Fatalf("old menu joined the new vote: %+v", result)
	}
	discordVotes.Lock()
	defer discordVotes.Unlock()
	if len(discordVotes.restore.voters) != 1 {
		t.Fatal("stale selection added a vote")
	}
}
