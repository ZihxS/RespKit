# RespKit v1.0.0 — Architecture

## Non-negotiable contracts

- **Import and use:** `resp.OK(ctx, data)` and `resp.Fail(ctx, err)` need no registration or setup.
- **Default schema:** envelope with `snake_case` RespKit-owned keys; application `data` keys are never rewritten.
- **Custom options:** naming, locale catalogs, error resolver mapping, RFC 9457, and local debug details are opt-in.
- **Security:** internal error strings cannot escape in ordinary response bodies; explicit `WithExposeInternalErrors(true)` is **unsafe** outside an isolated development environment.
- **Framework neutrality:** a single `resp` import and structural transport detection, no external framework imports in the root module. `net/http`, Chi, Gorilla/mux, and Beego v2 use `resp.HTTP(w, r)`.
- **Core compatibility:** Go 1.22+, standard-library-only root module.

## Dependency graph

```text
                                 +---------------------------+
     project handler -----------> | resp/ public facade        |
                                 | options / engine          |
                                 +---------------------------+
                                     |       |       |      |
                       +-------------+       |       |      +------------------+
                       v                     v       v                         v
                internal/transport    internal/resolver    internal/i18n    internal/formatter
                context bridge        semantic error       locale           JSON keys and
                + framework detection classification       matching        body templates
                                             |
                                        internal/diagnostics
                                          safe output for AI/CLI
```

The `internal` tree is an implementation detail: integrations and application code must use exported `resp` APIs only. Internal packages are allowed to import the standard library but never import `resp`, avoiding cycles.

## Runtime pipeline

1. Public `resp.OK`, `resp.Created`, `resp.NoContent`, `resp.Fail` (or methods on an immutable `Engine`).
2. `resp.Fail` uses typed application errors and `errors.Is`/`errors.As` in the resolver. Unknown errors safely map to a generic 500; ambiguous DB/cache errors are not blindly mapped.
3. Transport selects a framework bridge by structural interface assertion. Reflection is limited to framework-specific call shape differences.
4. `Accept-Language` plus optional per-response override select localized static catalog messages; client-facing code remains locale independent.
5. Formatter constructs the envelope or RFC 9457 representation; it **never** mutates the caller's `data` keys.
6. Body JSON is encoded before writing the HTTP status. For `204` and `HEAD`, omit body.
7. CLI diagnostics are metadata only, with no raw application error strings.

## Source layout

```text
resp/
  config.go       # immutable Engine instances and optional atomic default swap
  response.go     # public zero-config helpers and response options
  engine.go       # orchestrates the HTTP response pipeline
  errors.go       # public domain errors and wrapped causes
  resolve.go      # resolver facade and diagnostics
  i18n.go         # locale facade / response-level fallback
  transport.go    # net/http pairing, export compatibility
  catalog.go      # builtin error catalog
internal/
  transport/      # adapter.go, detector.go, gin.go, echo.go, fiber.go; stdlib wrappers cover Chi/mux/Beego
  resolver/       # typed errors, builtins, mapper ordering, SQLSTATE
  i18n/           # embedded per-locale JSON catalogs, matcher.go, translate.go
  formatter/      # envelope, response field naming, RFC 9457
  diagnostics/    # AI-/CLI-readable JSON
cmd/respkit/      # doctor, catalog, explain, version
examples/nethttp/ # runnable real HTTP example
integrations/     # separate Go modules for integration tests and optional adapters
  framework-tests/ # upstream runtime tests, organized by framework + shared testkit
  pgx/            # optional PostgreSQL error mapper
.github/workflows/# standard Go and real-framework CI
```

## Contributor extension rules

- Keep additions to HTTP framework support in `internal/transport/`; `resp` does not import external frameworks.
- Keep language catalogs and matching in `internal/i18n/`; store each built-in locale in its own embedded JSON file.
- Add typed error mapping in `internal/resolver/` only for deterministic, clearly safe classifications. For third-party-specific errors, prefer optional adapters or structural interfaces.
- Every mapping requires tests of wrapped errors, invalid input, locale fallback, and confidentiality.
- `resp.Configure` replaces the singleton atomically. Never mutate existing `Engine` maps/slices after creation or register mappings during requests.
- Pin framework versions in `integrations/framework-tests/go.mod` and verify *actual compilation and execution*; signature doubles are not proof of upstream compatibility.

## Known limits

- Pinned upstream framework tests are exercised from `integrations/framework-tests`; optional production adapters live beside them in their own modules. See `docs/VERIFICATION.md` for the versions actually run in this workspace.
- Fiber uses a small reflected `Status(int)` call; avoid claiming completely reflection-free dispatch.
- Framework contexts are single-request values, and response writing is a one-shot operation. RespKit cannot undo an already committed status.
- `100.0%` refers to statement coverage across the root `resp` package and `internal/...` during its regression tests, not full branch coverage, untested dependency versions, CLI commands, or external code.
