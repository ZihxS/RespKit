# Optional integrations

This directory groups integration-related modules under one name. The core module remains standard-library-only and Go 1.22-compatible. Each subfolder has its own `go.mod`, so test frameworks and optional adapters stay out of the core dependency graph.

## Available adapters

- [`pgx`](pgx/README.md): maps pgx/v5 unique-constraint (`23505`) and exclusion-constraint (`23P01`) errors to HTTP 409.

`framework-tests/` contains test-only end-to-end checks against upstream HTTP frameworks, organized one directory per framework/version. Other subfolders, such as `pgx/`, contain optional production adapters that applications can import. This keeps both test harnesses and adapters together without adding their dependencies to the core module.

The core already recognizes SQLSTATE `23505` through the structural `SQLState() string` interface. The pgx mapper adds explicit pgx diagnostics and maps `23P01`. Other database errors remain unmapped because their HTTP meaning depends on application-specific rules. The presence of this directory does not mean that GORM, validator, Redis, or other adapters have shipped.

Adapters require a clear source error type, documented dependency versions, opt-in enablement, no runtime global mutation, and tests for wrapped errors plus fallback behavior.
