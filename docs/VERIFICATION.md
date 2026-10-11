# RespKit verification report

This report records completed checks and the scope of each compatibility claim.

## Local verification

The latest recorded local run was performed on 11 October 2026 with Go 1.27.1 on macOS arm64. The core module supports Go 1.22 and later; its tests also passed locally with Go 1.22.12. The framework-test module requires Go 1.26 or later. The SQLite adapter supports Go 1.22 and requires cgo and a C compiler.

| Check | Result | Command |
| --- | --- | --- |
| Core tests | Pass | `make test` |
| Core tests with Go 1.22.12 | Pass | `GOTOOLCHAIN=go1.22.12 go test -count=1 ./...` |
| Race detector | Pass | `make race` |
| Static analysis | Pass | `make vet` |
| Combined statement coverage for `resp` and `internal/...` | 100.0% | `make coverage` |
| Linux/amd64 and Windows/amd64 cross-build smoke | Pass; compile only | `GOOS=linux/windows GOARCH=amd64 go test -exec=true -count=1 ./...` |
| GitHub Actions workflow syntax | Pass | `actionlint .github/workflows/ci.yml` |
| Upstream framework tests with race detector | Pass; 64 tests | `go test -race -count=1 -v ./...` in `integrations/framework-tests` |
| Optional pgx and SQLite mapper tests with race detector; SQLite vet | Pass; SQLite statement coverage 100.0% | `make integration-mappers` |
| Optional pgx static analysis | Pass | `go vet ./...` in `integrations/pgx` |
| SQLite mapper/runtime tests with Go 1.22.12 and race detector | Pass | `GOTOOLCHAIN=go1.22.12 CGO_ENABLED=1 go test -race -count=1 ./...` in `integrations/sqlite3` |
| Documentation build and internal links | Pass; 15 documentation pages checked | `npm run build` in `book` |

The SQLite runtime suite used `mattn/go-sqlite3 v1.14.52`. It executed real unique/primary-key/rowid, not-null, foreign-key, check, STRICT datatype, mismatch, busy/shared-cache lock, read-only, size/capacity, and database-open/invalid-file failures. HTTP assertions cover Indonesian/English messages, request locale negotiation, custom-language fallback, envelope/Problem Details formats, cancellation, and confidentiality. Synthetic errors cover I/O and out-of-memory conditions without disrupting the host environment.

The Go 1.22 SQLite run passed with a macOS linker `LC_DYSYMTAB` warning. The documentation build passed with two CSS optimizer warnings about the dependency's `::highlight` selector. Neither warning caused a failing check.

The cross-build smoke command compiles test binaries from macOS. It does not run them on Linux or Windows.

## Hosted CI

All hosted CI jobs use Go 1.27.x. The [quality-gates run on 11 October 2026](https://github.com/ZihxS/RespKit/actions/runs/38072670243) passed on Ubuntu, macOS, and Windows. It also passed the upstream framework tests, race detector, 100% statement-coverage gate, and optional pgx mapper tests.

An earlier workflow exposed a macOS runner issue with Go 1.22 and a coverage gap. CI now uses one Go version across all operating systems, and regression tests cover the previously unmeasured package paths. Go 1.22 remains the supported minimum for the core module and passed its local test run.

## Pinned upstream framework versions

| Framework or router | Tested version |
| --- | --- |
| Gin | v1.12.0 |
| Echo | v4.16.0 and v5.4.0 |
| Fiber | v2.52.14 and v3.3.0 |
| Chi | v5.3.2 |
| Gorilla/mux | v1.8.1 |
| Beego | v2.3.10 |

The suite sends HTTP requests through these routers and checks status codes, JSON bodies, localization and cache headers, error confidentiality, `Location`, and empty 204 responses. The result applies only to the pinned versions shown above.

Run the framework suite with `make integration`. It fetches upstream dependencies, so network access is required. Both pgx and SQLite mapper tests use separate modules. These modules do not add dependencies to the core package. The workflow is configured to test SQLite on Linux, macOS, and Windows with Go 1.27 and its Go 1.22 minimum on Linux; those newly added hosted jobs have not run as part of this local verification.

## Coverage scope

The core 100% figure covers combined statement coverage for the root `resp` package and runtime packages under `internal/`. It does not imply complete branch or path coverage. Framework tests and optional database mappers live in separate modules and are reported independently. SQLite's 100% statement coverage applies only to its mapper package, excluding upstream driver code.

## Security

The release review on 11 October 2026 found ten open Dependabot alerts affecting Fiber v2/v3 and quic-go in the isolated framework-test module. Its dependencies were updated to Fiber v2.52.14/v3.3.0 and quic-go v0.59.1, the fixed versions identified by those alerts. These framework dependencies are not part of the core or SQLite adapter's runtime dependency graph. All 64 framework tests and framework-module `go vet` passed with the updated versions. The shared test client now creates outgoing HTTP requests with `http.NewRequest`, preserving the Host header required by the updated fasthttp parser.

Raw error messages are hidden by default. `WithExposeInternalErrors(true)` adds a `debug_message` field and can disclose sensitive values; use it only in isolated local development. It is not enabled by request headers or inferred environment variables.
