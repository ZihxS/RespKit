# Changelog

## v1.0.0

- Kept the `resp.*` helpers usable without configuration.
- Organized framework transport, error resolution, localization, response formatting, and CLI diagnostics into internal packages.
- Separated Gin/Echo/Fiber transport bridges into their own source files.
- Kept `WithExposeInternalErrors(true)` as a deliberate local-development-only unsafe option (disabled by default).
- Retained default Indonesian/English localization, optional casing, deterministic error mapping, and optional RFC 9457.
- Enforced 100% statement coverage for `resp` plus all runtime `internal/...` packages through the combined regression suite.
- Added a separate integration module with tests against upstream frameworks and a CI workflow. The pinned framework tests and pgx mapper tests pass locally; see `docs/VERIFICATION.md` for details.
- Added architecture guide, contributor rules, local verification record and public release checklist.
