# S7UI product shell and API v3

Status: **Implemented on v7-nightly**

Recorded: **2026-08-08**

## Product shape

S7UI is one Svelte application with client-side paths. Navigation never reloads the document, but refresh, browser history and deep links retain the active view.

The shell is organized around user intent:

- Home
- Operate: Console, Logs and Activity
- Configure: Game, Gallery and Files
- Protect: Backups
- Access: People and access management
- System: Settings, Backends and experimental Plugins

Navigation and actions remain permission-aware. The active backend is workspace context rather than a hidden setting.

## Visual language

The interface uses layered glass for navigation, overlays and broad surfaces. Editors, consoles and dense data remain opaque enough to read for long sessions.

The curated presets are Självlysande, Skogsgrön, Norrskensdröm and Solkatt. Självlysande is the default. Custom image uploads affect Home only.

The old theme token collection is intentionally not migrated. Saved backend profiles and secure Electron credentials remain compatible.

## API boundary

'/api/v3' directly supersedes the earlier API routes. There is no v1/v2 compatibility contract on v7-nightly.

- Browser and Electron authentication use '/api/v3/auth'.
- Server, settings, runfile, gallery, file, backup, identity and plugin resources use '/api/v3'.
- SSE endpoints live below '/api/v3/streams' and send JSON event payloads.
- '/api/v3/overview' provides the permission-protected Home snapshot.
- Unknown '/api' paths return a JSON error instead of falling through to the SPA.
- Start, stop and restart are POST operations that return JSON.

The frontend deliberately keeps one API client for browser cookies and Electron's protected bearer-token bridge.

## Deliberate limits

Activity currently records frontend operations for the active application session. Durable backend task history can be added when long-running operations gain a common job model.

Plugins remain fully trusted native code. When the unsafe feature flag enables them, the UI presents an explicit warning; it does not imply sandboxing.

A backend continues to own one game server. Multi-server orchestration is not part of this package.
