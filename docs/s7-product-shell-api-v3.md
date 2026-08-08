# S7UI product shell and API v3

Status: **Implemented on v7-nightly; product validation pending**

Recorded: **2026-08-09**

## Product shape

S7UI is now one Svelte application with client-side routing, deep links and browser history. The active route survives refresh. Navigation and actions derive from the authenticated user's backend-issued permissions and feature flags.

The shell is organized by user intent:

- Home
- Operate: Console, Logs and durable Activity
- Configure: Game, Gallery and Files
- Protect: Backups
- Access: People and access management
- System: Settings, Backends and experimental Plugins

Every active screen uses the S7 material and interaction system. The unreachable pre-S7 application, auth, dashboard, gallery and settings components were removed so there is one frontend path to maintain.

## Visual language

Navigation, overlays and broad surfaces use layered glass. Editors, consoles and dense data remain opaque enough for long sessions.

The curated presets are Självlysande, Skogsgrön, Norrskensdröm and Solkatt. Självlysande is the default. Custom background uploads affect Home only. Saved backend profiles and secure Electron credentials remain compatible.

## Product bootstrap

After authentication the shell retrieves `GET /api/v3/capabilities`. It declares:

- the API version and backend build version;
- the one-server-per-backend boundary;
- enabled product features;
- supported cross-cutting capabilities;
- the current principal's ordered permission set.

Connection health uses this permission-neutral authenticated endpoint. A user without `server.view` is no longer shown as disconnected merely because server status is forbidden.

## API v3 boundary

`/api/v3` directly supersedes earlier API routes. There is no v1/v2 compatibility contract on v7-nightly.

Regular v3 handlers are normalized to JSON. Errors use `{"error":{"code","message"}}`; resource results use a `data` envelope. SSE endpoints live below `/api/v3/streams` and carry JSON event payloads.

Important resource contracts:

| Resource | Contract |
| --- | --- |
| `/api/v3/auth/*` | browser sessions, desktop credentials, users, groups, tokens, sessions and audit |
| `/api/v3/capabilities` | authenticated product and permission manifest |
| `/api/v3/overview` | Home snapshot |
| `/api/v3/activity` | capped durable operation history |
| `/api/v3/streams/activity` | live operation events |
| `/api/v3/settings` | typed GET catalog and PATCH mutation with effect metadata |
| `/api/v3/files` | editable-file catalog |
| `/api/v3/files/content` | bounded GET and PUT textual content |
| `/api/v3/backup/*` | status, validated create/restore and stable restore-point metadata |
| `/api/v3/streams/*` | console, event and log streams |

Unknown `/api` paths return a JSON error instead of falling through to the SPA. State-changing endpoints remain protected by session CSRF/origin validation or a scoped bearer credential.

## Activity model

Successful authenticated mutations at the v3 boundary create an actor-attributed operational event. The most recent 500 events are stored in `SSUI/config/activity-v3.json`, loaded after restart and broadcast to connected browser and Electron clients. The frontend merges that history with transient local notices and toasts.

This is operation history, not yet a general background-job engine. A request accepted for asynchronous work records acceptance; detailed progress and terminal job state remain future work.

## Browser and Electron

The renderer uses one API client. Browser deployments use secure session cookies and CSRF. Remote Electron profiles use bearer credentials kept in the main process, protected OS storage, explicit HTTPS-only backend URLs and per-origin certificate fingerprint approval. SSE uses the same main-process credential boundary.

The actual packaged entrypoint is `main-v7.cjs`; the legacy `main.cjs` is not selected by `package.json` or the builder manifest.

## Deliberate limits

Plugins remain fully trusted native code. They are disabled by default. When the unsafe feature flag enables them, the UI presents an explicit warning; this is not sandboxing.

A backend continues to own one game server. Multi-server orchestration is not part of this package.

## Validation completed

- Full Go test suite
- v3 JSON boundary tests
- durable activity write/reload/subscriber tests
- backup metadata, protected-file and setting-effect tests
- frontend production build
- Electron main and preload syntax checks
- tracked frontend audit for removed v1/v2 and superseded file/settings calls

Packaged Windows/Linux Electron smoke tests and visual workflow review still require release artifacts and human testing.
