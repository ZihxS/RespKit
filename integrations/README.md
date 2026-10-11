# Optional integrations

The core RespKit module supports Go 1.22 and later and uses only the standard library. Modules under this directory add framework test dependencies or optional production adapters without expanding the core dependency graph.

## Available modules

- [pgx](pgx/README.md) is an optional adapter that maps PostgreSQL errors from pgx/v5.
- [sqlite3](sqlite3/README.md) is an optional Go 1.22-compatible adapter for `mattn/go-sqlite3`, with real SQLite runtime tests and localized public messages. It requires cgo and a C compiler.
- `framework-tests` runs HTTP requests through real upstream framework routers. It is test-only and is not imported by applications.

The core module already recognizes PostgreSQL SQLSTATE `23505` through the `SQLState() string` interface. The pgx adapter adds explicit mappings for additional PostgreSQL errors and identifies its mapper in diagnostics. Other database errors remain unmapped unless an application enables an integration with a deterministic mapping.

This directory does not imply support for adapters that have not shipped. Every production adapter should use typed errors, remain opt-in, document its dependency and version requirements, and test wrapped errors and fallback behavior.

## Contribute an integration

See the [optional integrations roadmap](https://github.com/ZihxS/RespKit/issues/4) for integration tasks and their acceptance criteria. Follow the [contribution guide](../CONTRIBUTING.md) when submitting a pull request.
