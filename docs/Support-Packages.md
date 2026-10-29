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

## When SSUI does not start

Run SSUI with `-Diagnose`, for example `./StationeersServerUI -Diagnose`. It checks the installation, the config, the folders and their permissions, DNS, and runs the Spacecat connectivity check, prints the findings and the full data, and exits without starting SSUI or changing any file. The output is plain text, so you can paste it into your ticket or give it to an AI assistant together with [ssui.dev](https://ssui.dev). Known passwords and tokens are redacted.

## Open a support connection

This is how the SSUI team supports you. Expand **Support & Diagnostics**, pick a duration and click **Open support connection**. SSUI shows a code. In your Discord ticket, press **Start remote connection** and enter the code. While the connection is open, SSUI sends the team a snapshot every minute: server status, diagnostics, the redacted configuration, connected players with their Steam IDs, recent backups, mods, routines, the last log lines and system details: OS, kernel, Docker, the user SSUI runs as, and the owner, mount, filesystem and free space of the SSUI, save and log folders. The team can also request support packages, a connectivity check and a log recording, and start, stop or restart the game server. Passwords and tokens are never sent. Nothing else: no config changes, no console commands, no shell. The relay on jxsn.dev only passes the data through and does not keep it. Everything the team or the support bot fetches is posted in your ticket.

SSUI only opens an outgoing HTTPS connection to `jxsn.dev`, no port is opened. Every request appears in the panel while it is running. **End connection** stops it immediately, and it ends by itself after the chosen duration. Opening a connection requires `server.control`.

### The support bot during a remote connection

Our Discord support bot is an AI teammate. It knows the SSUI source code and documentation and answers in your ticket by default. While a remote connection is active, it also sees the live overview of your SSUI, the same data the support team sees, and uses it in its answers.

- **It can fetch data itself.** When it needs more information, the bot can request read-only data from your SSUI: status, diagnostics, report, logs, connectivity check, config (without passwords and tokens), players, backups, mods, routines and the game server state. Every result is posted in your ticket, marked "fetched by the support bot", so you always see what it looked at.
- **It can only suggest a server start, stop or restart.** The bot asks you in the ticket first, and only you, the person who opened the ticket, can confirm with **Yes, go ahead**. If you say **No**, or nobody answers within 10 minutes, nothing changes. The same applies when a supporter asks for it.

The bot never changes your configuration, files or saves, and it cannot run anything outside the support commands SSUI allows.

### Bug reports to GitHub

If a ticket turns up a real bug, the support team can ask the bot to draft a GitHub issue. The bot writes the draft in English, with a summary, steps to reproduce, expected and actual behavior, versions and relevant log excerpts, and posts it in the ticket. Only after a team member confirms it with **Create issue** is it filed in the matching repository, and the bot posts the link. Issue drafts never contain Discord names or IDs, IP addresses, Steam IDs or other personal data. Users cannot create issues through the bot.

### Privacy

The AI behind the bot is provided by OpenAI. What the bot sees in your ticket, the conversation and, with a remote connection, your SSUI status and the results it fetched, is sent there so it can answer you. IP addresses and Steam IDs are removed first; player names and log lines stay in. OpenAI's own retention rules apply there. The bot keeps the conversation of your ticket, including that data, a short summary of the ticket and the Discord ID of whoever started the connection for 180 days after the last activity in the ticket so it can follow up, then deletes them automatically. Only anonymous lessons learned, a problem and its fix without names, IDs, IP addresses or server names, are kept longer. The bot runs isolated from the rest of the SSUI infrastructure.

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
| `LockedBoxCat` | A SSUI, save or log folder is not writable or sits on a read-only mount. |
| `FullBowlCat` | Less than 2 GB are free for a SSUI, save or log folder. |
| `StrayCat` | A folder SSUI cannot use belongs to a different user. |
| `ScratchedCat` | SSUI cannot read or write some files or folders in its data or save folders. |
| `LostCat` | SSUI cannot resolve the host names it needs, for example for updates, SteamCMD or Discord. |

---

[Documentation home](index.md) · [All guides](index.md#run-your-server)
