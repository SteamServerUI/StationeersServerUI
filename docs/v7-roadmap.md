# S7UI v7 production roadmap

Status: **Planned**  
Target: **Christmas 2026**

## Product boundary

S7UI v7 is a reliable, approachable manager for one Steam game server per backend. Runfiles provide the game abstraction. A future orchestration layer may control multiple backends, but one backend managing multiple game processes is not part of v7.

The existing implementation is a substantial base, not a throwaway prototype. Server management, SteamCMD, runfiles, settings, backups, logs, Discord, TLS and the Svelte UI already exist. The work is to make their contracts secure, understandable and dependable.

## Phase 1: identity and trust

Release blockers:

- Replace the mixed JWT/cookie frontend authentication with a coherent browser-session protocol.
- Add backend-enforced permissions and custom access groups.
- Give Electron a credential flow designed for remote and eventually multiple backends.
- Add revocable, named and scoped API tokens.
- Require authentication and provide a local owner recovery command.
- Disable the plugin system by default until its trust contract is explicit.
- Establish integration tests around every protected boundary.

This phase deliberately comes before a visual redesign. The UI needs stable session and permission APIs to build on.

## Phase 2: application shell

- Replace duplicated frontend auth state with one session and backend client.
- Separate connection setup, certificate trust and login.
- Add routing, deep links and dependable browser history.
- Standardize loading, empty, error, success and confirmation states.
- Make navigation and actions permission-aware while keeping the backend authoritative.
- Define a small component and styling system before redesigning individual screens.

## Phase 3: core workflows

Improve the workflows in this order:

1. First-run owner setup
2. Runfile and game selection
3. Dashboard and server lifecycle
4. Console and logs
5. Game and backend settings
6. Backups and restore
7. Users, groups and tokens
8. Updates
9. Files and advanced configuration

Each workflow should be usable without knowing how the Go packages or runfile JSON are arranged internally.

## Phase 4: hardening and release

- Test upgrades and configuration migrations on Windows, Linux and Docker.
- Exercise interrupted updates, failed starts, corrupt runfiles and failed restores.
- Add browser and Electron end-to-end coverage for the critical workflows.
- Review secrets, logs, local sockets, uploads, file access and debug services.
- Run a closed beta with written entry and exit criteria.
- Ship release candidates only when every network endpoint has an owner, permission and test.

## Deferred work

- Multiple game servers inside one backend
- A central identity provider owned by SSUI
- Safe execution of untrusted native plugins
- Fully customizable dashboard layouts
- A general automation engine

These ideas should influence interfaces where doing so is cheap, but they must not delay a dependable v7.
