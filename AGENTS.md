# RespKit coding agent instructions

RespKit has a **Go 1.22-compatible, standard-library-only** root module. Its owner-selected module path is `github.com/ZihxS/RespKit`; keep it consistent across imports and nested modules.

## Public API invariants

- `resp.OK`, `resp.Created`, `resp.NoContent`, `resp.Fail`, `resp.NotFound` work without configuration.
- Gin/Echo/Fiber use their native contexts; Chi, Gorilla/mux, Beego, and `net/http` use `resp.HTTP(w,r)`.
- Default JSON keys use `snake_case`; only RespKit-owned keys change with the naming option.
- Default locale `id`; built-ins `id/en`; `Accept-Language` is parsed at request time.
- Public error codes are stable, language independent, and use HTTP status codes correctly.
- Raw error messages **MUST NOT** be exposed by default. `WithExposeInternalErrors(true)` is explicitly unsafe and reserved for isolated local development.

## Editing guidance

- Change adapters in `internal/transport` and validate with **actual** framework runtime tests.
- Change error classification in `internal/resolver`; use `errors.Is/As`; do not classify by error message strings.
- Change translation lookup in `internal/i18n`, body formatting in `internal/formatter`, machine-readable JSON in `internal/diagnostics`.
- Keep the `resp` public facade thin; do not couple it to external framework package imports.
- Do not add optional external service calls, secrets, or credentials to the package runtime.
- Maintain 100% combined statement coverage of `resp` plus `internal/...` measured via `make coverage`.

## Verify

```bash
make test
make race
make vet
make coverage
make integration  # Go 1.26+ and internet required; uses real Gin/Echo/Fiber/Chi/Gorilla mux/Beego
```

Do not assert upstream integration success unless the genuine integration tests actually ran. For future releases, follow `book/content/docs/release.mdx` and publish only a commit that passes every required CI check.
