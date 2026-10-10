# RespKit verification record

## Current workspace run

Run date: 2026-10-10. Environment: Go 1.27.1, macOS/darwin arm64. The root module remains compatible with Go 1.22; the framework-test module under `integrations/framework-tests` requires Go 1.26+.

All hosted CI jobs use Go 1.27.x. The root module still declares Go 1.22 compatibility, and its tests passed locally on Go 1.22.12. CI runs root tests and vet across Ubuntu, macOS, and Windows, plus root race/coverage checks and upstream framework and pgx mapper tests.

Local Linux/amd64 and Windows/amd64 smoke commands compiled all package and test binaries from the macOS/arm64 workspace. The `-exec=true` flag means those cross-platform commands did not execute the test binaries; hosted CI must confirm runtime behavior on each operating system.

| Check | Result | Command |
|---|---|---|
| Root tests | PASS | `make test` |
| Root tests on declared minimum Go version | PASS — Go 1.22.12 | `GOTOOLCHAIN=go1.22.12 go test -count=1 ./...` |
| Root race detector | PASS | `make race` |
| Root static analysis | PASS | `make vet` |
| Root runtime statement coverage | **100.0%** | `make coverage` |
| Linux/amd64 and Windows/amd64 cross-build smoke | PASS — compile-only | `GOOS=linux/windows GOARCH=amd64 go test -exec=true -count=1 ./...` |
| GitHub Actions workflow syntax | PASS | `actionlint .github/workflows/ci.yml` |
| Upstream framework tests with race detector | **PASS — 64 tests across 8 framework packages** | `go test -race -count=1 -v ./...` in `integrations/framework-tests` |
| Optional pgx mapper with race detector | **PASS** | `make integration-mappers` |

The integration suite exercises actual framework routers and verifies status codes, response JSON, localization headers, error confidentiality, `Location`, and empty 204 responses.

| Framework/router | Pinned version |
|---|---:|
| Gin | v1.12.0 |
| Echo | v4.16.0, v5.4.0 |
| Fiber | v2.52.12, v3.0.0 |
| Chi | v5.3.2 |
| Gorilla/mux | v1.8.1 |
| Beego | v2.3.10 |

Run the same suite with `make integration`; it first tidies the isolated module and then executes the race-enabled tests. Network access is required to resolve upstream modules. These dependencies remain outside the root module. An empty test config file prevents Beego's upstream package initialization from warning about a missing `conf/app.conf`.

## Core coverage gate

Coverage is measured over the root `resp` plus `internal/...` runtime and passed the 100% combined statement coverage gate in this workspace. The framework runtime tests and optional pgx mapper live in separate modules under `integrations/`. Statement coverage does not imply full branch or path coverage.

## Security note

Raw errors remain hidden by default. `WithExposeInternalErrors(true)` adds `debug_message`, can disclose secrets, and must only be used in private local development. It is not enabled by request headers or inferred environment variables.
