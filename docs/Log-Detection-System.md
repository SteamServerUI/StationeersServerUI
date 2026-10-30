# Logs and detections

**Console** shows Stationeers output. **Events** turns recognized lines into saves, connections and other readable activity. **Logs** shows SSUI's own subsystem messages. Start with the stream for the component that failed.

Built-in detections also drive connected-player and server-state information. If Stationeers changes its log format, that interpretation can lag behind the actual game. Check the raw console before assuming an empty player list is accurate.

## Add a custom event

Open **Config → Detection Manager**. Each entry is a routine: one trigger and up to ten actions that run from top to bottom. Managing routines requires `detections.manage`.

1. Copy a real console line that identifies the condition, or pick a built-in server event.
2. Choose the trigger: **Log contains** (case-sensitive literal), **Log matches regex** (captures variable values) or **Server event** (player joined/left, world saved, server joinable and others).
3. Add actions in the order they should run, then save.
4. Trigger the condition and check Events, Discord and the game.

| Action | What it does |
|---|---|
| Post to event log | Shows the text in Events and the Discord event log channel |
| Post to Discord hub | Posts the text in the Discord hub channel |
| Announce to players | In-game `announce` (needs SSCM and a running server) |
| Run console command | Any SSCM console command, for example `kick {1}` |
| Save world | SSCM `save` |
| Wait | Pauses the chain for 1-3600 seconds |
| Call webhook | POSTs `{"content", "routine", "values"}` as JSON; Discord webhook URLs work as-is |
| Start, stop, restart server | Same sequence as the Discord hub buttons |
| Restore 2nd-newest backup | Stops the server, restores the second-newest backup and starts it again |

If an action fails, the rest of the chain is skipped and the failure is posted to the event log. A webhook counts as failed when it does not answer with a 2xx status within 10 seconds or its HTTPS certificate is not trusted.

| Trigger | Pattern | Action text |
|---|---|---|
| Log contains | `Unsupported shader` | `Unity reported an unsupported shader.` |
| Log matches regex | `Player (.+) has reached level (\d+)` | `Player {1} reached level {2}` |
| Server event | Player joined | `Welcome aboard, {1}!` |

The regex example is illustrative log text, not a claim that Stationeers has that event. For `Player Alex has reached level 12`, it produces `Player Alex reached level 12`.

Use Go-compatible regex syntax. Numbered capture groups become `{1}`, `{2}` and so on, and `{0}` is the whole match. For player events, `{1}` is the player name and `{2}` the SteamID. Test both matching and unrelated lines. `.*` will match impressively often and tell you almost nothing.

A routine ignores new triggers until its cooldown has passed. Routines that start, stop, restart or restore the server always wait at least two minutes, and they share one slot with the Discord hub, so two server actions never overlap.

Definitions are stored at `SSUI/config/customdetections.json`. Include it when moving your configuration. Detections from older versions become routines with a single event-log action. Add/remove changes apply while SSUI is running.

## Diagnose a detection

| Symptom | Check |
|---|---|
| Nothing fires | Exact case, spaces and punctuation; confirm the line exists in Console |
| Regex rejected | Regex syntax; API JSON needs an additional level of backslash escaping |
| Missing capture text | Pattern defines and matches the requested group |
| Too many events | Narrow the pattern to a distinctive fragment |
| UI event but no Discord event | Discord enabled, event channel ID and bot permissions |

For a broken built-in detection, report the raw line, SSUI version and game build. A screenshot of the last error in a cascade is less useful than the first error in text. [Collect diagnostics](Support-Packages.md).

## API and implementation

`/api/v3/streams/events` is the authenticated SSE stream. Custom definitions use GET/POST/DELETE `/api/v3/detections`. Request details are in [API](API.md#settings-and-maintenance).

The implementation lives in `src/managers/detectionmgr`: keyword and regex handlers parse console text, publish structured events and update relevant runtime state. Raw console output remains available independently of whether a pattern matches.

---

[Documentation home](index.md) · [All guides](index.md#run-your-server)
