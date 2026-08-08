# S7UI v7 development notes

This folder holds the plans and decisions for turning the current v7 alpha into a production release. The documents are meant to survive individual implementation branches and keep the reasons behind larger changes visible.

## Documents

- [v7-roadmap.md](v7-roadmap.md) describes the path from the current beta-era implementation to a stable v7 release.
- [identity-sessions-authorization.md](identity-sessions-authorization.md) is the implementation contract for users, access groups, sessions, API tokens and Electron authentication.
- [security-considerations.md](security-considerations.md) tracks security boundaries that need deliberate treatment before a stable release.
- [implementation-status.md](implementation-status.md) records what the first identity and trust implementation actually shipped and what validation remains.

## Status words

- **Idea**: worth exploring, but not accepted work.
- **Planned**: accepted direction; details may still change.
- **In progress**: implementation has started.
- **Implemented**: present in v7-nightly and covered by tests.
- **Release blocker**: stable v7 must not ship without resolving it.

Update the documents when a decision changes. Do not quietly let the code and plan disagree.
