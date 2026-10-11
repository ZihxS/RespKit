# RespKit coding agent instructions

RespKit has a **Go 1.22-compatible, standard-library-only** root module. Its owner-selected module path is `github.com/ZihxS/RespKit`; keep it consistent across imports and nested modules.

## Go idioms and code style

- Write idiomatic Go compatible with each module's declared minimum version. Run `gofmt` on changed Go files and `go vet` in every changed Go module.
- Prefer simple control flow, early returns, explicit types, and small functions. Introduce helpers or named types when they remove meaningful duplication; avoid speculative abstractions and dependencies.
- Use clear Go names and document exported identifiers. Use standard-library constants for HTTP statuses instead of unexplained numeric literals in new code.
- Use `resp.Code...` constants for built-in codes in public definitions, application errors, optional mappers, and Go examples. Internal runtime code uses the canonical constants in `internal/resolver/codes.go`. Define application-specific codes and protocol values as descriptive constants when the dependency provides no named value; avoid magic strings and numbers in classifications. Keep wire values stable and keep translation JSON keys aligned with them.
- For composite literals with multiple elements or fields, put **one element or field per line**, including test tables, slices, maps, and nested struct literals. Use keyed struct fields rather than positional multi-field literals. Separate fields in struct declarations onto their own lines as well.
- Use slices for ordered collections and test cases. Use `map[T]struct{}` as a set only when uniqueness or membership lookup is required; map iteration order is unspecified. Prefer deterministic test order unless varying order is the purpose of the test.
- Keep fixed test groups separate and iterate them directly when they share an expectation. Do not build a combined table with `append` or repeat an anonymous struct definition just to add a static group. `append` remains idiomatic for collections that actually grow; preallocate when the size is known and doing so makes the code clearer or measurably helps.
- Preserve error identity with `%w` and classify using `errors.Is`/`errors.As`. Handle errors explicitly and keep internal error text out of public messages.
- Respect context cancellation and resource ownership. Register test resource cleanup with `t.Cleanup`; avoid leaking database connections, transactions, goroutines, or mutable global state.
- Keep engine configuration immutable after creation. Add concurrency mechanisms only when shared state requires them, and verify concurrency changes with the race detector.
- Optimize based on measurements and actual requirements. Prefer clear switches and direct loops for small fixed collections; do not replace them with maps or generic frameworks without a demonstrated benefit.

Apply these conventions to new and edited code without unrelated repository-wide formatting changes.

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
- Keep driver-specific error mappers in separate opt-in modules under `integrations/`, using the public `resp` API. Test actual driver runtime errors, wrapped errors, fallback behavior, localization, and confidentiality.
- Do not add optional external service calls, secrets, or credentials to the package runtime.
- Maintain 100% combined statement coverage of `resp` plus `internal/...` measured via `make coverage`.

## Verify

```bash
make test
make race
make vet
make coverage
make integration  # Go 1.26+ and internet required; uses real Gin/Echo/Fiber/Chi/Gorilla mux/Beego
make integration-mappers  # Optional pgx/SQLite modules; SQLite requires cgo and a C compiler
```

Do not assert upstream integration success unless the genuine integration tests actually ran. For future releases, follow `book/content/docs/release.mdx` and publish only a commit that passes every required CI check.
