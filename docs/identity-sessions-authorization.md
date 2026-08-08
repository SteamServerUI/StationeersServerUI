# Identity, sessions and authorization

Status: **In progress**

## Goals

- Know which human or service is making every protected request.
- Enforce permissions in the backend, not just by hiding frontend controls.
- Support custom access groups without complicated deny precedence.
- Give browsers safe sessions and orchestration clients safe scoped tokens.
- Keep identities local to a backend while leaving room for a future multi-backend controller.
- Make credentials revocable and privileged changes auditable.

## Identity storage

Identity data lives in the main atomically-written configuration. The identity section has its own schema version so migrations do not depend on the rest of the settings format.

Persist:

- Users with UUID, normalized unique username, display username, Argon2id password hash, enabled state, group IDs and timestamps
- Groups with UUID, name, description, system flag and permission grants
- Browser sessions as hashes with user ID, CSRF hash, creation, last use, idle expiry and absolute expiry
- Named tokens as identifier plus secret hash, owner, scopes, creation, expiry, last use and revocation time

Never persist raw passwords, raw session IDs, raw CSRF values or raw token secrets. Sensitive comparisons must be constant-time. Identity mutations use the existing configuration lock and atomic save path.

The built-in Owner group is immutable and grants every permission. Other groups are custom. A user may belong to several groups and receives the union of their grants. There are no deny grants.

## Permission model

Permissions describe backend actions rather than UI pages. The initial catalog covers:

- Server status, lifecycle and console input
- Logs and connected-player information
- Backup viewing, creation, restore and deletion
- Runfile viewing and management
- Settings viewing and management
- File reading and writing
- SteamCMD and SSCM execution
- Plugin viewing and management
- User viewing and management
- Group viewing and management
- Own-token and all-token management
- Session management
- Audit-log viewing
- Backend reload, update and security administration

Every protected route declares one permission. Permission middleware runs after authentication and places the resolved principal and effective permissions in the request context. A missing credential returns `401`; an authenticated principal without the grant returns `403`.

## First setup and migration

Existing alpha and beta users, JWTs and API keys are not migrated. On the first start with the new identity schema:

1. Legacy authentication material becomes inactive.
2. The backend enters restricted setup mode.
3. A short-lived one-time setup secret is generated and printed locally.
4. Only health and owner-bootstrap endpoints are available.
5. Bootstrap consumes the secret, creates the first Owner and invalidates the setup secret.

Authentication cannot be skipped. Lost ownership is recovered through an audited local CLI command which invalidates active credentials before resetting or creating an owner.

## Browser sessions

Browser login creates a random opaque session ID. The raw ID is stored only in an `HttpOnly`, `Secure`, `SameSite=Strict` cookie. The backend stores its SHA-256 hash.

Sessions have idle and absolute expiration, survive backend restarts and can be revoked independently. Password changes, account disablement and owner recovery revoke affected sessions. Logout is idempotent.

Successful login and session inspection return user details, group summaries, effective permissions, expiration and a CSRF value. The CSRF value lives only in frontend memory and is required in a custom header for state-changing cookie-authenticated requests. The backend also validates the request origin.

The browser never receives a session ID in JSON and stores no credential in local storage. Backend URLs and non-secret display preferences may remain there.

## Named tokens

Electron, automation and future orchestration use bearer tokens rather than browser cookies. A token belongs to a user or service identity and has a name, explicit scopes and optional expiry. Its scopes cannot exceed the creator's permissions.

The token format contains a public lookup identifier and a random secret. Only the secret hash is stored. The complete token is displayed once. Tokens record last use and support independent revocation and rotation. Browser-session endpoints do not accept query-string tokens.

## Electron

Electron treats each backend as a separate connection profile. The renderer does not receive bearer tokens.

- The main process encrypts credentials with Electron `safeStorage` and stores them under the application data directory.
- A context-isolated preload API exposes a narrow request and authentication bridge.
- The main process performs authenticated backend requests and returns sanitized responses.
- Certificate validation stays enabled. Self-signed backends require an explicit fingerprint trust flow per backend.
- The same profile structure can later serve a multi-backend orchestration view without changing backend authentication.

## HTTP contract

- `/auth/login` creates a browser session and returns session metadata.
- `/auth/logout` revokes the current browser session.
- `/api/v2/auth/session` reports the current principal, groups, permissions, CSRF value and expiry.
- User, group, session and token resources use permission-protected `/api/v2/auth/...` endpoints.
- Responses use a consistent JSON error object.
- Handlers reject unsupported methods.
- CORS uses an explicit origin policy and never reflects arbitrary credentialed origins.

## Acceptance criteria

- No protected endpoint lacks a declared permission and integration test.
- Raw credentials never appear in config, logs, URLs or renderer storage.
- Restarted backends retain valid sessions; expired and revoked sessions fail.
- Password reset and disabled accounts invalidate their credentials.
- Custom group changes affect authorization immediately.
- Browser, bearer and setup credentials cannot be substituted for one another.
- Electron supports several saved backend profiles without disabling TLS validation.
