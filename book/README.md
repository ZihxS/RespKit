# RespKit

RespKit adalah pustaka Go untuk membentuk response HTTP yang konsisten, menerjemahkan pesan error, dan mengubah error Go menjadi status HTTP yang aman.

```go
resp.OK(ctx, data)
resp.Created(ctx, data, resp.Location("/users/123"))
resp.NoContent(ctx)
resp.Fail(ctx, err)
resp.NotFound(ctx)
```

Satu API ini dapat dipakai dengan `net/http`, Gin, Echo v4/v5, dan Fiber v2/v3. Chi, Gorilla/mux, Beego v2, serta handler `net/http` memakai adapter `resp.HTTP(w, r)`.

## Yang menjadi default

- Tidak perlu membuat instance atau mendaftarkan framework untuk pemakaian dasar.
- Response berhasil berbentuk `{"success":true,"data":...}`.
- Kode error stabil dan tidak bergantung pada bahasa.
- Bahasa default `id`; katalog bawaan tersedia untuk Bahasa Indonesia dan Inggris.
- Pesan error internal tidak dikirim ke client.
- Modul utama menggunakan Go 1.22+ dan standard library saja.

## Mulai membaca

- [Memulai](memulai.md)
- [Response berhasil dan gagal](response.md)
- [Pemetaan error](errors.md)
- [Bahasa dan terjemahan](localization.md)
- [Konfigurasi](configuration.md)
- [HTTP framework](frameworks.md)
- [PostgreSQL dan pgx](postgres.md)
- [CLI diagnostics](diagnostics.md)
- [Keamanan](security.md)
- [Testing dan kontribusi](testing.md)

> **Status distribusi:** dokumentasi ini menjelaskan API RespKit v1.0.0. Repository source tersedia di [GitHub](https://github.com/ZihxS/RespKit), tetapi belum ada release tag. Instalasi tetap dapat dilakukan menggunakan versi Go module terbaru.
