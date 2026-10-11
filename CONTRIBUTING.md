# Contributing to RespKit

RespKit provides consistent HTTP responses and error mapping for Go applications. The core
module supports Go 1.22 and later and has no third-party runtime dependencies. Keep framework
and database dependencies in separate modules under integrations/.

The [optional integrations roadmap](https://github.com/ZihxS/RespKit/issues/4) lists open
tasks for databases, libraries, and HTTP frameworks. Each task includes its scope and
acceptance criteria. The contributor-guide task is labeled good first issue.

## Project areas

- `internal/transport`: framework detection and adapters. Gin, Echo, and Fiber use native
  contexts; net/http, Chi, Gorilla/mux, and Beego use the explicit `resp.HTTP(w, r)` adapter.
- `internal/resolver`: typed error classification and stable public codes. Use
  `errors.Is` and `errors.As`; do not infer status codes from error text.
- `internal/i18n`: locale matching and JSON translation catalogs. Do not translate raw Go
  error messages automatically.
- `internal/formatter`: response shapes, key naming, and RFC 9457 Problem Details. Preserve
  application data as provided.
- `internal/diagnostics`: structured error metadata without secrets.
- `resp`: public helpers, engine options, and compatibility contracts. Review exported API
  changes for backward compatibility.
- `book/content/docs`: user and contributor documentation published on GitHub Pages.

## Before opening a pull request

The main branch requires pull requests and passing CI checks, including for administrators.
Keep the branch up to date with main before merging. Force pushes and deletion of main are
disabled.

1. Include regression tests for changed behavior and real upstream tests for framework
   adapters.
2. Use typed errors and stable, language-independent codes. Keep raw errors out of client
   responses by default.
3. Keep third-party dependencies in a separate module under integrations/.
4. Update relevant documentation in English when behavior, setup, or compatibility changes.
5. Preserve the existing exported API or explain the breaking change clearly.

## Run the checks

Run from the repository root:

```bash
make test
make race
make vet
make coverage
make integration
make integration-mappers
```

The coverage gate requires 100% combined statement coverage for resp and internal runtime
packages. The framework test module requires Go 1.26 or later; the optional pgx adapter
requires Go 1.25. The optional SQLite adapter supports Go 1.22 and requires cgo and a C
compiler. `make integration-mappers` checks both database adapters. Integration checks
download upstream dependencies and need network access.

Build the documentation site from the book directory:

```bash
npm ci
npm run build
```

The documentation build checks internal links before generating the static site. See
[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) for package boundaries and
[docs/VERIFICATION.md](docs/VERIFICATION.md) for tested versions and results.

CI also runs integration tests against the actual upstream framework routers. Local interface
doubles help test edge cases, but they do not replace the upstream suite.
