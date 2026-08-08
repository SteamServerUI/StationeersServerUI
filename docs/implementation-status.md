# S7UI v7 implementation status

Status: **Implemented on v7-nightly**

Recorded: **2026-08-08**

The identity, browser session, authorization and desktop credential foundation is now active. This is the base for the remaining v7 product work, not a statement that all of v7 is production-ready.

## Present now

- Local users with Argon2id password hashes, immutable Owner access and additive custom groups
- Backend permission enforcement across the v7 HTTP API
- Persistent opaque browser sessions with secure cookies, idle and absolute expiry, CSRF and origin checks
- Named bearer tokens with explicit scopes, expiry, last-use tracking and revocation
- First-run Owner bootstrap with a short-lived local setup secret
- Local Owner recovery that revokes every session and token
- Persistent capped audit history for setup, recovery, user, group, token and session changes
- Permission-aware People & access UI for users, groups, tokens, sessions and audit events
- Electron profiles whose bearer credentials stay in the main process and require protected OS storage
- Explicit per-backend certificate fingerprint approval without a global TLS bypass
- Plugins, their management endpoints and their local registration socket disabled by default

Owner recovery reads `SSUI_RECOVERY_PASSWORD` only when `-RecoverOwner` (or `-r`) is supplied, then removes the variable from the process environment. The password is never accepted as a command-line value or written to logs.

Plugins can currently be restored only by setting `SSUI_ENABLE_UNSAFE_PLUGINS=true`. This means fully trusted native code with the same operating-system access as SSUI; it is not a sandbox or a production plugin security model.

## Validation still required for a stable release

- Browser and packaged Electron end-to-end tests on Windows and Linux
- Upgrade, migration and rollback exercises using real v5 and early-v7 configurations
- Complete method, size-limit and traversal review for legacy file, upload and command handlers
- Windows named-pipe ACL design before the unsafe plugin system can be enabled outside isolated development
- Packaging, updater and certificate-trust smoke tests on release artifacts

The one-game-server-per-backend boundary remains unchanged. Multi-server orchestration is still deferred.
