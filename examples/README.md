# Examples

Run `go run ./examples/nethttp` from the repository root, then try:

```bash
curl -i -H 'Accept-Language: en-US' http://localhost:18090/not-found
curl -i -H 'Accept-Language: id-ID' http://localhost:18090/internal
```

Gin, Echo, Fiber, Chi, Gorilla/mux, and Beego are covered by executable HTTP tests in `integrations/framework-tests`. Shared response assertions are in `integrations/framework-tests/internal/testkit`. These tests use a **separate Go module**, so importing the root library does not pull in framework dependencies.
