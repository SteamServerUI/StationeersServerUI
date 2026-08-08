# S7UI v7 implementation status

Status: **Foundation and S7 product shell implemented on v7-nightly**

Recorded: **2026-08-09**

This ledger says what exists in code. It does not call v7 production-ready before packaged, migration and destructive-workflow validation.

## Identity and trust

- Local users with Argon2id password hashes, immutable Owner access and additive custom groups
- Shared browser/desktop login throttling with uniform-cost checks for unknown and disabled users
- Backend permission enforcement across the v7 HTTP API
- Persistent opaque browser sessions with secure cookies, idle and absolute expiry, CSRF and origin checks
- Named bearer tokens with explicit scopes, expiry, last-use tracking and revocation
- First-run Owner bootstrap with a short-lived local setup secret
- Local Owner recovery that revokes every session and token
- Persistent capped security audit history
- Permission-aware People & access UI
- Electron credentials kept in the main process and protected OS storage
- Per-backend HTTPS certificate fingerprint approval without a global TLS bypass
- Plugins, management endpoints and local registration socket disabled by default

Owner recovery reads `SSUI_RECOVERY_PASSWORD` only when `-RecoverOwner` or `-r` is supplied, then removes it from the process environment.

## S7 product experience

- One routed application shell with permission-aware navigation and deep links
- Modern Home, Console, Logs, Activity, Game, Gallery, Files, Backups, People & access, Settings, Backends and Plugins views
- Four curated themes and Home-specific custom backgrounds
- One browser/Electron API client
- Capability-backed backend identity and connection health
- Persistent operational activity history plus live SSE delivery
- Typed settings with explicit application/reload effects
- Bounded JSON file editor resource
- Restore-point metadata and validated backup actions
- Unreachable pre-S7 application and component trees removed

## API v3

- v3 is the only frontend API contract; it supersedes v1 and v2
- Regular v3 responses are JSON-normalized
- Stream payloads are JSON SSE
- Successful protected mutations are actor-attributed in operational activity
- Unknown API paths return structured JSON errors
- Plugins advertise as a disabled feature unless the unsafe opt-in is set

## Validation completed locally

- `go test ./...`
- Frontend `npm run build`
- `node --check main-v7.cjs`
- `node --check preload.cjs`
- Contract tests for the JSON boundary, activity persistence, backup formats, protected files and setting effects
- Source audit showing no active v1/v2 frontend calls

## Still required before stable

- Browser and packaged Electron end-to-end tests on Windows and Linux
- Upgrade, migration and rollback exercises using real v5 and early-v7 configurations
- Real-game smoke tests for start, stop, restart, SteamCMD, runfile changes and log streams
- Destructive restore tests including interrupted and corrupt backups
- Complete upload, command and file-path hardening review
- Windows named-pipe ACL design before unsafe plugins can be enabled outside isolated development
- Packaging, updater and certificate-trust smoke tests on release artifacts
- Closed-beta visual and workflow review at common desktop sizes

Plugins can currently be restored only with `SSUI_ENABLE_UNSAFE_PLUGINS=true`. They execute as fully trusted native code with SSUI's operating-system access.

The one-game-server-per-backend boundary remains unchanged.
