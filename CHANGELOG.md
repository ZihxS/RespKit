# Changelog

## v1.1.0 — 2026-10-11

[GitHub release](https://github.com/ZihxS/RespKit/releases/tag/v1.1.0)

- Added public `resp.Code...` constants for built-in response and validation codes. Core classification, localization fallback, and optional database mappers reuse the canonical values; HTTP definitions use named `net/http` statuses.
- Added an opt-in `integrations/sqlite3` module for `mattn/go-sqlite3`, compatible with Go 1.22. Typed and wrapped errors map to stable public codes with Indonesian/English messages and support for custom locales; raw database errors remain hidden by default.
- Added real SQLite runtime tests, integration documentation, and CI checks for the SQLite mapper.
- Aligned the mattn SQLite mapper with issue #8: I/O failures and unknown extended codes stay safe internal errors, and the integration documentation explains driver choice.
- Published SQLite and pgx as separately versioned optional modules at `integrations/sqlite3/v1.0.0` and `integrations/pgx/v1.0.0`, requiring core `v1.1.0`. The core remains Go 1.22-compatible and standard-library-only; SQLite requires cgo and a C compiler, while pgx requires Go 1.25.
- Updated Go coding guidance to require named error/status constants, keyed composite literals with one value per line, and deterministic test tables.
- Updated the isolated framework-test dependencies to Fiber `v2.52.14` / `v3.3.0` and quic-go `v0.59.1` to address the Dependabot alerts reviewed for this release.

## v1.0.0

[GitHub release](https://github.com/ZihxS/RespKit/releases/tag/v1.0.0)

- Kept the `resp.*` helpers usable without configuration.
- Organized framework transport, error resolution, localization, response formatting, and CLI diagnostics into internal packages.
- Separated Gin/Echo/Fiber transport bridges into their own source files.
- Kept `WithExposeInternalErrors(true)` as a deliberate local-development-only unsafe option (disabled by default).
- Retained default Indonesian/English localization, optional casing, deterministic error mapping, and optional RFC 9457.
- Enforced 100% statement coverage for `resp` plus all runtime `internal/...` packages through the combined regression suite.
- Added a separate integration module with tests against upstream frameworks and a CI workflow. The pinned framework tests and pgx mapper tests pass locally; see `docs/VERIFICATION.md` for details.
- Added architecture and contributor guides, a verification report, and release documentation.
