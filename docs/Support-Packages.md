# Support packages

A support package is a ZIP with a readable `report.html` and a machine-readable `diagnosis.json`. It shows the SSUI version, branch, OS, game build, a redacted view of your settings and every finding that needs attention. A plain package contains **no logs**.

## Create one in the WebUI

Open **Config**, scroll to the bottom and expand **Support & Diagnostics**. The findings are shown right there; anything not listed passed or does not apply.

- **No logs needed:** click **Download support package**.
- **Logs needed:** click **Start recording**, reproduce the issue, click **Stop recording**, then **Download recording**.

A recording captures every SSUI line, including debug lines, and all gameserver output, no matter what your normal log settings are. It stops capturing at 64 MB, and SSUI keeps the last three recordings in `SSUI/logs/support-recordings`. The diagnostics in a recording are taken when you stop it.

Creating packages and recordings requires the `server.control` permission.

## Create one in SSCLI

These are terminal commands, not Stationeers console commands.

1. Enter `supportmode` (`sm`) to start a recording.
2. Reproduce the issue.
3. Enter `supportpackage` (`sp`). It stops the recording and writes the ZIP to the folder SSUI runs from.

`sp` without a recording writes a package without logs.

## Open a support connection

If the SSUI team asks for it in your Discord ticket, expand **Support & Diagnostics**, pick a duration and click **Open support connection**. SSUI shows a code; enter it in your ticket with `/ssui connect`. While the connection is open, SSUI sends the team a snapshot every minute: server status, diagnostics, the redacted configuration, connected players with their Steam IDs, recent backups, mods, routines and the last log lines. The team can also request support packages, a connectivity check and a log recording. Passwords and tokens are never sent. Nothing else: no remote control, no config changes, no shell. The data is relayed to the ticket and not stored.

SSUI only opens an outgoing HTTPS connection to `jxsn.dev`, no port is opened. Every request appears in the panel while it is running. **End connection** stops it immediately, and it ends by itself after the chosen duration. Opening a connection requires `server.control`.

## Before you share it

Include what you did and when the first error appeared. Use the [issue tracker](https://github.com/SteamServerUI/StationeersServerUI/issues) or [Discord](https://discord.gg/8n3vN92MyJ).

!!! warning
    SSUI replaces the server password, auth secret, admin password, Discord token and JWT key in recorded logs with `[REDACTED]`. Player names, Steam IDs, paths, IP addresses and anything else a subsystem printed stay in. Open the ZIP and review it before sharing it.

## Finding codes

| Code | Meaning |
|---|---|
| `RootCat` | SSUI runs as root outside a container. |
| `HairballConfig` | The config file is missing or unreadable. |
| `EmptyBoxCat` | The Stationeers executable or its game data is missing. |
| `TangledYarnCat` | The world settings are incomplete, for example a start condition without a difficulty. |
| `LostWorldCat` | The world settings are not in the installed game's catalog, or the catalog could not be read. |
| `Folder404Cat` | The SSUI data folder is missing or unreadable. |
| `Log4Cat` | File logging is on, but the log folder is unusable. |
| `SteamedCat` | SteamCMD is missing. |
| `CatnipMissing` | SSCM is enabled, but its plugin folder is missing. |
| `SaveMeow` | The save name is invalid, or the save folder does not exist yet. |
| `Portcat` | A port or the local IP is invalid, or two ports collide. |
| `Spacecat` | The last connectivity check could not reach the server from outside. |
| `Jittercat` | The last connectivity check reached only some UDP packet sizes. |
| `SchrodingersCat` | There is no connectivity check result yet. |

---

[Documentation home](index.md) · [All guides](index.md#run-your-server)
