# Contributing to RespKit

RespKit is a small library for HTTP responses and error mapping. Keep third-party framework dependencies out of the root module.

## Where to contribute

- `internal/transport`: HTTP framework detection and bridges (Gin, Echo, Fiber, net/http). Chi, Gorilla/mux, and Beego use the explicit `resp.HTTP(w,r)` adapter. Add upstream framework tests in `integrations/framework-tests`.
- `internal/resolver`: conservative typed error mapping and safe error codes. Domain-specific errors must be explicit; ambiguous errors should fall back safely.
- `internal/i18n`: language matching, catalogs and translations. No automatic translation of raw Go errors.
- `internal/formatter`: response shape, casing, HTTP Problem Details; never mutate application `data` fields.
- `internal/diagnostics`: machine-readable diagnostics with no secrets.
- `resp`: public API, options, facade and stable contracts. Changes here require compatibility review.

## Requirements for a pull request

1. Include source-library tests and any relevant real-framework integration tests.
2. Run `make test`, `make race`, `make vet`, and `make coverage` (100% runtime statement coverage gate).
3. Use `errors.Is`/`errors.As`, never parse raw `err.Error()` strings to guess status codes.
4. Include stable error codes, translations, correct HTTP statuses, and no raw error leaks by default.
5. If you add a third-party dependency, isolate it in a module under `integrations/`, away from the root module.
6. Preserve existing exported API or explicitly explain the breaking change.
7. Update `docs/ARCHITECTURE.md` and `docs/VERIFICATION.md` when changing components or test claims.

CI runs a separate integration job against upstream framework routers. Tests using local interface doubles do not replace that job.
