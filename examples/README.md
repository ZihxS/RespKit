# Examples

Run `go run ./examples/nethttp` from the repository root, then try:

```bash
curl -i -H 'Accept-Language: en-US' http://localhost:18090/not-found
curl -i -H 'Accept-Language: id-ID' http://localhost:18090/internal
```

The genuine Gin/Echo/Fiber, Chi, Gorilla/mux, and Beego examples are executable end-to-end HTTP tests in their own `integrations/framework-tests` subdirectories. Shared response assertions live in `integrations/framework-tests/internal/testkit`. Their dependencies live in a **separate Go module**, so the root library remains standard-library-only and importable without pulling framework code.
