# Public v1.0.0 release gate

- [x] Separate public API from transport, localization, error resolver, formatter, diagnostics.
- [x] Preserve default `resp.OK`, `resp.Fail`, and zero-config API.
- [x] Keep public `data` fields independent from casing convention.
- [x] Keep raw errors hidden by default; explicit unsafe local debugging option available.
- [x] Pass Go unit, race, vet, and 100% statement coverage gate for `resp` + `internal/...`.
- [x] Execute race-enabled real integration tests against upstream Gin, Echo v4/v5, Fiber v2/v3, Chi, Gorilla/mux, and Beego v2 with Go 1.26+; see [`VERIFICATION.md`](VERIFICATION.md) for pinned versions and results.
- [x] Select module path `github.com/ZihxS/RespKit` and use it throughout the repository.
- [x] Configure CI with Go 1.27.x across Ubuntu, macOS, and Windows, plus race/coverage, upstream framework tests, and pgx mapper tests.
- [x] Confirm the latest GitHub Actions run passes on Ubuntu, macOS, and Windows ([run 38072373586](https://github.com/ZihxS/RespKit/actions/runs/38072373586)).
- [x] Add maintainer security contact `m.saleh.solahudin@gmail.com` to [`SECURITY.md`](../SECURITY.md).
- [ ] Enable GitHub Private Vulnerability Reporting and Dependabot alerts, then protect the default branch with CI checks.
- [ ] Create GitHub release tag `v1.0.0` only after checks above succeed.

The repository contains the v1.0.0 source, but no GitHub release has been published. A passing CI run covers only the pinned versions listed in `VERIFICATION.md`.
