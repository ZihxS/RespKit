# Public v1.0.0 release gate

- [x] Separate public API from transport, localization, error resolver, formatter, diagnostics.
- [x] Preserve default `resp.OK`, `resp.Fail`, and zero-config API.
- [x] Keep public `data` fields independent from casing convention.
- [x] Keep raw errors hidden by default; explicit unsafe local debugging option available.
- [x] Pass Go unit, race, vet, and 100% statement coverage gate for `resp` + `internal/...`.
- [x] Execute race-enabled real integration tests against upstream Gin, Echo v4/v5, Fiber v2/v3, Chi, Gorilla/mux, and Beego v2 with Go 1.26+; see [`VERIFICATION.md`](VERIFICATION.md) for pinned versions and results.
- [x] Select module path `github.com/ZihxS/RespKit` and use it throughout the repository.
- [x] Configure CI matrix: root tests and vet on Go 1.22/1.27 for Ubuntu, macOS, and Windows; race/coverage on Linux; framework and pgx tests on their minimum Go versions.
- [ ] Confirm the hosted GitHub Actions matrix passes after this source is pushed to its repository.
- [x] Add maintainer security contact `m.saleh.solahudin@gmail.com` to [`SECURITY.md`](../SECURITY.md).
- [ ] After the repository exists, enable GitHub Private Vulnerability Reporting and Dependabot alerts, then protect the default branch by requiring pull requests and CI checks.
- [ ] Create GitHub release tag `v1.0.0` only after checks above succeed.

The packaged source is labeled v1.0.0 for requested API/versioning; it is **not** a published or upstream-certified release.
