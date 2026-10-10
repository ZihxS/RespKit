# Changelog

## v1.0.0 — refactored source candidate

- Preserved the zero-configuration `resp.*` public surface.
- Extracted framework transport, error resolution, localization, response formatting, and AI-readable diagnostics into narrowly focused internal packages.
- Separated Gin/Echo/Fiber transport bridges into their own source files.
- Kept `WithExposeInternalErrors(true)` as a deliberate local-development-only unsafe option (disabled by default).
- Retained default Indonesian/English localization, optional casing, deterministic error mapping, and optional RFC 9457.
- Enforced 100% statement coverage for `resp` plus all runtime `internal/...` packages through the combined regression suite.
- Included optional real framework test module and CI gate. **The upstream suite has not been executed in this offline environment.**
- Added architecture guide, contributor rules, local verification record and public release checklist.
