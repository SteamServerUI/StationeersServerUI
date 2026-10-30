# v6 changelog

SSUI v6 is the biggest update StationeersServerUI has had so far. The core stays the same: game server management, configuration and the loader work as they did in v5. Around that core, a lot has grown. The backup manager now polls for autosaves and never misses one, the Discord bot became a server hub with community votes, SSUI has real user accounts with permission groups, and custom detections grew into routines that can announce, call webhooks and even restart your server. On top of that come a connectivity check for your game port, built-in diagnostics, an opt-in support connection to the SSUI team, a cleaner API, new backgrounds and a fully localized config page.

If you only read one thing before upgrading, read the [upgrade guide](Upgrade-Guide.md): the API moved to v3, runtime data moved from `UIMod` to `SSUI` and Docker now uses a single mount. Everything below is the full developer changelog from v5.14.1 to v6.0.0, grouped by area.

---

## Breaking changes

- break(runtime): move managed runtime data from `UIMod/` to `SSUI/`; copy it once on first start and leave `UIMod/` untouched for rollback
- break(api): replace `/api/v2` with `/api/v3`; all responses use a shared `data`/`error` envelope
- break(api): state-changing endpoints require non-GET methods
- break(api): remove the old v2 routes (`/saveconfig`, `/saveconfigasjson`, `/start`, `/stop`, `/console`, `/events`, `/logs/*`, `/api/v2/*`)
- break(backups): identify backups by filename; numeric backup indexes are rejected instead of guessed
- break(discord): `/restore` and `/download` take `name:<filename>` instead of `index`
- break(discord): remove the reaction control panel and `controlChannelID`/`controlPanelChannelID`; the status panel channel becomes the hub
- break(discord): admin actions require the configured `discordAdminRoleID`; Discord's Administrator permission no longer bypasses it
- break(auth): replace legacy sessions, JWTs and API keys with users, groups, sessions and scoped personal access tokens
- break(backups): replace the retention keys (`backupKeepLastN`, `backupKeepDailyFor`, `backupKeepWeeklyFor`, `backupKeepMonthlyFor`, `backupCleanupInterval`, `isCleanupEnabled`, `backupWaitTime`); old keys are not migrated and cleanup starts disabled
- break(docker): one writable `./server-data:/app:rw` bind mount is the supported layout; split mounts are not migrated
- break(cli): remove the terminal dashboard and `IsCLIDashboardEnabled`
- break(module): move the Go module to `github.com/SteamServerUI/StationeersServerUI/v6`

---

## Runtime and platform

- refactor(runtime): rename the runtime folder to `SSUI` and migrate `UIMod` once with an incomplete-migration marker
- refactor(runtime): skip symlinks during migration instead of following them
- feat(runtime): handle SIGTERM and interrupts; stop the game and backup manager within 25 seconds, second signal forces exit
- fix(runtime): harden runtime storage and config writes with temp file, sync and rename
- fix(runtime): replace files atomically on Windows and Unix through platform-specific helpers
- refactor(config): make config updates transactional through `PATCH /api/v3/settings`
- feat(config): add `LoadConfigWithoutSaving` for read-only tooling
- fix(loader): detect cached `go run` builds in the Go build cache as development workdirs
- fix(loader): coalesce Discord bot reloads so concurrent reloads do not race
- fix(loader): take the bot's application ID from the session before registering slash commands
- refactor(loader): run the startup sanity check through the diagnostics package
- fix(loader): detect Docker even when `-NoSanityCheck` is set
- feat(loader): log failing and warning diagnostics findings after startup
- fix(loader): fix the "exiting in 10 seconds" message
- build(go): pin the go directive to 1.26.0 so older toolchains can fetch it

## Identity, users and access

- feat(auth): add the v6 identity store at `SSUI/security` with users, groups, sessions, tokens and an audit trail
- feat(auth): hash passwords with Argon2id (64 MiB, 3 iterations); keep verifying old bcrypt hashes
- feat(auth): import v5 users into the identity store on first start and record the migration in the audit log
- feat(auth): add 20 permissions (`server.view`, `server.control`, `console.read`, `console.write`, `backups.view`, `backups.analyze`, `backups.download`, `backups.restore`, `settings.view`, `settings.manage`, `steamcmd.run`, `slp.manage`, `detections.manage`, `update.install`, `users.manage`, `groups.manage`, `tokens.manage`, `audit.view`, `backend.reload`, `security.manage`)
- feat(auth): add the protected `system-owner` group
- feat(auth): add sessions with 24 hour idle and 30 day absolute lifetime and CSRF tokens derived per session
- feat(auth): add scoped personal access tokens (`ssui_pat_...`) with optional expiry
- feat(auth): rate-limit logins to 5 failures per 15 minutes per user and address, then block for 15 minutes
- feat(auth): validate usernames (3 to 64 characters, letters, digits, `-`, `_`, `.`) and require 10 character passwords
- feat(auth): compare against a dummy hash for unknown users to avoid timing leaks
- feat(auth): add `-RecoveryPassword` recovery owner that is removed again when the flag is gone
- refactor(setup): simplify initial owner setup to a single bootstrap step
- fix(auth): persist identity state reliably on Windows
- fix(auth): send empty user groups, group permissions and token scopes as `[]` instead of `null`
- fix(web): serve only login and setup assets publicly; everything else requires a session
- feat(web): add People & Access in the config page with account, people, groups and API token panels
- feat(web): show only the pages and controls the signed-in user may use
- feat(web): rework the expanded user editor and fix missing spacing between access cards
- fix(web): sign out, token revoke and other People & Access actions threw on every click because the tab container carried `data-delete-user`

## API v3

- feat(api): add `GET /api/v3` with a shared response contract
- feat(api): add auth routes: `login`, `logout`, `session`, `sessions`, `password`, `setup`, `setup/bootstrap`, `users`, `groups`, `tokens`, `audit`
- feat(api): add `GET/PATCH /api/v3/settings`
- feat(api): add server routes: `status`, `players`, `start`, `stop`, `connectivity`, `connectivity/check`
- feat(api): add backup routes: `backups`, `backups/analysis`, `backups/download`, `backups/restore`
- feat(api): add detection routes: `GET/POST /api/v3/detections`, `DELETE /api/v3/detections/{id}`
- feat(api): add SLP routes: `slp/mods`, `slp/mods/update`, `slp/packages`, `slp/install`, `slp/reinstall`, `slp/uninstall`
- feat(api): add `sscm/status`, `sscm/commands`, `steamcmd/run`, `backend/reload`, `advertiser/override`, `tls/certificate`, `worldgen/catalog`, `update`
- feat(api): add support routes: `support/diagnostics`, `support/package`, `support/recording`, `support/recording/start|stop|package`, `support/connection`
- feat(api): move log and console streams to `/api/v3/streams/{console,events,logs/*}`
- feat(web): add `GET /healthz`
- refactor(web): move the WebUI and Svelte client to the v3 API
- docs(api): expand the API reference with request examples

## Backup Manager v4

- feat(backup): replace fsnotify with polling of the autosave directory every 45 seconds
- feat(backup): handle a file only after two unchanged observations at least 45 seconds apart
- feat(backup): validate `world_meta.xml` and `world.xml` structure and ZIP checksums before and after copying
- feat(backup): copy into a temp file and publish under the final name only when complete
- feat(backup): persist backup metadata in `backup-meta-manifest.ssui` and keep a RAM inventory
- feat(backup): rebuild a missing manifest; keep a malformed one as `.invalid-...`
- feat(backup): add cached deep save analysis (days, players, things, rooms, atmospheres, pipe and cable networks, furnaces)
- feat(backup): stream world XML during analysis instead of holding it in memory
- feat(backup): backfill missing or outdated analyses in the background, newest first, batched every two seconds
- feat(backup): prioritize new saves over the analysis backlog
- feat(backup): recopy an archive that disappears while its original still exists
- feat(backup): retry unavailable storage instead of treating it as empty
- feat(backup): carry pending observations and retention state across reloads when paths stay the same
- feat(backup): replace retention with `backupRetentionEnabled`, `backupKeepNewestCount`, `backupDailyRetentionDays`, `backupWeeklyRetentionWeeks`, `backupMonthlyRetentionMonths`, `backupCleanupIntervalHours`
- feat(backup): apply retention rules as a union on local calendar days, ISO weeks and months
- feat(backup): delete source autosaves older than 24 hours only when an analyzed safe copy exists
- fix(backup): block cleanup on storage or manifest errors
- fix(backup): reject negative retention values and restore safe defaults
- fix(backup): protect backups from stale inventory
- fix(backup): reject overlapping source and archive folders, including through symlinks
- fix(backup): check restores before stopping the server and again under the backup lock
- fix(backup): only hold the backup lock while publishing, so restores no longer wait for validation, copy and analysis
- feat(backup): log progress at info and debug level under the `BACKUP` subsystem
- feat(backup): limit concurrent downloads
- feat(web): add the dedicated Backups page at `/backups` with expandable statistics and icons
- feat(web): show the latest backup on the dashboard
- test(backup): add rotation, reload, crash recovery, retention, manifest and selection tests

## Discord

- feat(discord): replace the control panel with one server hub message that refreshes every 15 seconds
- feat(discord): show state, players, version, uptime, next restart and statistics from the latest backup in the hub
- feat(discord): move save statistics to their own hub button
- feat(discord): add private Server Admin Actions: start, stop, restart, update, backups, announce, console command, ban, unban
- feat(discord): confirm stop, restart, update and restore
- feat(discord): add `/restart`
- feat(discord): add a private backup browser with ten files per page and single-use selection tokens
- feat(discord): reuse the existing hub message and search the last 100 messages for it
- feat(discord): share one execution slot for admin actions and votes; reject conflicts instead of queueing
- feat(discord): add restart and restore votes with duration, threshold, minimum and cooldown settings
- feat(discord): open a temporary public panel per vote and remove it 15 minutes after it ends
- feat(discord): restore votes offer the three newest backups and restart after success
- feat(discord): add backup statistics emojis
- feat(discord): post admin actions with the Discord user ID to the event channel
- feat(discord): keep the 10 MiB download cap and reject oversized files before reading them
- fix(discord): repeat role checks on confirmations and form submissions
- style(discord): improve information density in the hub
- style(discord): replace em dashes with plain hyphens
- feat(discord): expose `RunServerAction` for routines and the support connection
- test(discord): add hub, vote and backup selection tests

## Detection Manager routines

- feat(detections): turn custom detections into routines with one trigger and up to ten actions
- feat(detections): add triggers: log contains, log matches regex, built-in server event
- feat(detections): add actions: event log, Discord hub, announce, console command, save, wait, webhook, start, stop, restart, restore second-newest backup
- feat(detections): pass regex groups and player name/SteamID as `{0}`, `{1}`, `{2}`
- feat(detections): add cooldowns; server actions wait at least two minutes and share the Discord slot
- feat(detections): skip the rest of a chain and report it when an action fails
- feat(detections): treat webhooks as failed without a 2xx answer within 10 seconds or with an untrusted certificate
- feat(detections): migrate old detections to routines with a single event-log action
- feat(web): rebuild the Detection Manager UI
- fix(detections): custom detections failed to save with the v3 API
- fix(detections): refresh the list after saving instead of calling the removed `showTab`
- fix(detector): detect session registration across game builds and capitalizations
- fix(detector): treat a completed world save as proof of a healthy server
- fix(detector): match current and historical save-completion messages

## Spacecat connectivity check

- feat(connectivity): add Spacecat, an external UDP probe of the gameserver port
- feat(connectivity): probe with 1000, 1200 and 1400 byte packets to spot MTU problems
- feat(connectivity): run the check while the gameserver is stopped and reserve the port during it
- feat(connectivity): add `ConnectivityCheckEnabled` (default on)
- feat(web): show Spacecat on the dashboard with "Run check again"
- refactor(network): move public IP resolution into a shared package used by the advertiser and Spacecat
- docs(connectivity): note that the probe runs in Frankfurt and its latency says nothing about player speed

## Support and diagnostics

- feat(diagnostics): add one diagnostics system for the startup check, API, support and `-Diagnose`
- feat(diagnostics): report findings with code, status and fix, e.g. `EmptyBoxCat`, `SaveMeow`, `Spacecat`
- feat(diagnostics): collect OS, kernel, container signals, UID/GID, CPUs and RAM
- feat(diagnostics): inspect working, SSUI, saves and log folders: mount, filesystem, read-only flag, free space, owner and mode
- feat(diagnostics): add `LockedBoxCat` (folder not writable or read-only), `FullBowlCat` (under 2 GB free) and `StrayCat` (owner mismatch)
- feat(diagnostics): scan file access in data and save folders (`ScratchedCat`)
- feat(diagnostics): resolve GitHub, Steam CDN, jxsn.dev and Discord (`LostCat`)
- feat(diagnostics): add `-Diagnose` to print a redacted text report, with exit code 1 on failure
- feat(diagnostics): add support packages with `diagnosis.json` and `report.html`
- feat(diagnostics): add log recordings that include SSUI and gameserver logs in the package
- refactor(diagnostics): one package format; logs only via explicit recording
- feat(diagnostics): redact known passwords and tokens
- refactor(cli): route `supportmode`/`supportpackage` to the recorder
- feat(support): add an opt-in, time-limited support connection to jxsn.dev (30 minutes to 2 hours)
- feat(support): pair the connection with a Discord ticket through a one-time code
- feat(support): push a redacted snapshot every minute and answer read-only commands
- feat(support): let the support team start, stop and restart the gameserver with the ticket opener's consent
- feat(support): show requests from the support team in the UI
- fix(support): close quietly when the team ends the connection
- fix(web): show "waiting" while the connection registers instead of "not reachable"
- feat(web): move Support & Diagnostics into its own config tab
- docs(support): describe the support bot, issue drafts, privacy and 180 day retention

## WebUI

- style(web): unify the config page visual identity
- feat(web): separate expert settings behind a switch
- feat(web): discover worlds, difficulties, start conditions and locations from game files, with a bundled fallback
- feat(web): stabilize dashboard controls and resource links
- feat(web): add a page loader
- refactor(web): rebuild the SSUI update dialog
- feat(web): stop config tabs from widening on hover; show the name as a tooltip
- feat(web): add a phone icon for support and a workflow icon for the Detection Manager
- feat(web): localize all config page sections, the update dialog and expert settings in en-US, de-DE and sv-SE
- feat(background): replace the old background with a WebGL celestial view that shows the configured world
- feat(background): add station panoramas, classic planets and no background
- feat(background): add reduced motion and fall back to it without hardware WebGL
- feat(background): pause rendering when the window loses focus
- feat(background): rebuild the background menu with tabs and switches next to its button
- feat(seasonal): replace `xmas.js` with a shared seasonal system for Halloween and Christmas
- feat(seasonal): show a Halloween and a Christmas planet in the celestial view
- feat(seasonal): show the small decorations only with other backgrounds
- feat(seasonal): add a 30 day opt-out cookie that only applies to the current season
- feat(seasonal): run Halloween from October 20 to November 5
- docs(credits): add recent donors

## StationeersLaunchPad and workshop

- feat(slp): overhaul workshop mod management
- feat(slp): accept several workshop handles or URLs at once and drop duplicates
- feat(slp): copy only items whose SteamCMD download succeeded
- fix(slp): report missing SteamCMD and failed or missing items clearly

## Updater

- refactor(update): rebuild the SSUI update flow
- feat(update): list stable, major and prerelease candidates separately
- feat(update): require separate confirmation for major and prerelease updates
- feat(update): download the exact release that was reviewed
- feat(update): report image updates in Docker instead of replacing the binary

## Build, release and Docker

- build(release): make release builds tag driven and create assets from version tags
- build(release): support prerelease versions
- build(docker): prepare container and release distribution
- build(docker): copy the binary into `/app` atomically in the entrypoint
- build(docker): use one `/app` mount in Compose with a 30 second stop grace period
- ci(docs): publish the documentation from nightly

## CLI

- refactor(cli): rename `ssuicli` to `sscli`
- refactor(cli): remove the terminal dashboard (`dashboard` command)

## Documentation

- docs: move the documentation to ssui.dev
- docs: add 71 pages and assets, including upgrade, Docker, API, backup, Discord, Detection Manager, support and security guides
- docs: keep the v5 pages under Legacy
- docs: point README at ssui.dev

## Dependencies

- deps: update discordgo to v0.29.0
- deps: use `golang.org/x/crypto/argon2` for password hashing
- deps: remove fsnotify, golang-jwt and the charmbracelet TUI stack

---

[Documentation home](index.md) · [Upgrade from v5](Upgrade-Guide.md) · [Downloads](Downloads.md)
