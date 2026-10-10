# RespKit verification report

This report records completed checks and the scope of each compatibility claim.

## Local verification

The latest recorded local run was performed on 10 October 2026 with Go 1.27.1 on macOS arm64. The core module supports Go 1.22 and later; its tests also passed locally with Go 1.22.12. The framework-test module requires Go 1.26 or later.

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
| Optional pgx mapper tests with race detector | Pass | `make integration-mappers` |

The cross-build smoke command compiles test binaries from macOS. It does not run them on Linux or Windows.

## Hosted CI

All hosted CI jobs use Go 1.27.x. The [quality-gates run on 11 October 2026](https://github.com/ZihxS/RespKit/actions/runs/38072670243) passed on Ubuntu, macOS, and Windows. It also passed the upstream framework tests, race detector, 100% statement-coverage gate, and optional pgx mapper tests.

An earlier workflow exposed a macOS runner issue with Go 1.22 and a coverage gap. CI now uses one Go version across all operating systems, and regression tests cover the previously unmeasured package paths. Go 1.22 remains the supported minimum for the core module and passed its local test run.

## Pinned upstream framework versions

| Framework or router | Tested version |
| --- | --- |
| Gin | v1.12.0 |
| Echo | v4.16.0 and v5.4.0 |
| Fiber | v2.52.12 and v3.0.0 |
| Chi | v5.3.2 |
| Gorilla/mux | v1.8.1 |
| Beego | v2.3.10 |

The suite sends HTTP requests through these routers and checks status codes, JSON bodies, localization and cache headers, error confidentiality, `Location`, and empty 204 responses. The result applies only to the pinned versions shown above.

Run the framework suite with `make integration`. It fetches upstream dependencies, so network access is required. The pgx mapper tests also use a separate module. Neither module adds dependencies to the core package.

## Coverage scope

The 100% figure covers combined statement coverage for the root `resp` package and runtime packages under `internal/`. It does not imply complete branch or path coverage. Framework tests and the optional pgx mapper live in separate modules and are reported independently.

## Security

Raw error messages are hidden by default. `WithExposeInternalErrors(true)` adds a `debug_message` field and can disclose sensitive values; use it only in isolated local development. It is not enabled by request headers or inferred environment variables.
