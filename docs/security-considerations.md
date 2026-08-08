# S7UI v7 security considerations

Status: **Planned**

This is a living threat register, not a claim that every item is resolved.

## Release blockers

### Plugins

The current plugin gallery downloads executables from URLs in a remote manifest, marks them executable and runs them as the SSUI user. Artifacts have no signature or pinned digest. Plugin UI is proxied into an unsandboxed same-origin iframe, request credentials are forwarded, and plugins can reach the privileged local socket API.

Until that contract is redesigned, plugins must be disabled by default and described as fully trusted native code. A production design needs verified manifests and artifact hashes, atomic downloads, explicit permissions, credential-stripping proxies, isolated UI origins and a capability-limited backend API. OS process isolation is desirable but is not promised for the first stable v7.

### Authentication and authorization

The beta implementation mixes cookies, frontend JWT storage, bearer headers and query-string tokens. Authorization is all-or-nothing. Replace it with the protocol in [identity-sessions-authorization.md](identity-sessions-authorization.md).

### TLS and origins

- Remove Electron's global certificate-error bypass.
- Do not reflect arbitrary origins while allowing credentials.
- Keep browser sessions same-origin.
- Make self-signed certificate trust an explicit fingerprint decision.
- Add HTTP read, header, write and idle timeouts.

## Important boundaries

### Local sockets and named pipes

The local API exposes powerful backend operations. Filesystem permissions on Linux and an explicit named-pipe security descriptor on Windows must restrict access. A native plugin running as the SSUI user remains fully trusted unless a narrower plugin protocol replaces this API.

### Files and uploads

File read, save and upload handlers need canonical-path enforcement, size limits, content rules where appropriate and tests for traversal and symlink escapes. TLS key and certificate uploads need stricter permissions than ordinary runfile assets.

### Commands and game arguments

SteamCMD, SSCM, runfile arguments and server controls cross process-execution boundaries. Avoid shell construction, validate structured inputs and restrict these endpoints with precise permissions.

### Secrets and logs

Passwords, cookies, bearer tokens, setup secrets, Discord tokens, TLS private keys and game credentials must not appear in normal or debug logs. Configuration exports and support bundles need redaction.

### Debug services

Debug-only HTTP routes, plugin registration and pprof must not bind publicly by default. Starting a debug service should be obvious in logs and configuration.

## Release review checklist

- Inventory every HTTP, SSE, socket and pipe endpoint.
- Assign authentication type, permission and method constraints.
- Verify request-size limits and server timeouts.
- Test path traversal, malformed JSON and oversized inputs.
- Check browser CSP, iframe sandboxing and injected remote content.
- Review dependencies and update/build provenance.
- Exercise setup, recovery, token revocation and account-disable paths.
- Verify production builds cannot silently enable debug surfaces.
