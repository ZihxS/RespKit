# Konfigurasi

Pemakaian paling sederhana tidak memerlukan konfigurasi. Helper `resp.OK`, `resp.Fail`, dan lainnya memakai engine default.

## Engine khusus

Gunakan `resp.New` untuk memisahkan kebijakan antar API atau menghindari konfigurasi global:

```go
api, err := resp.New(
	resp.WithNaming(resp.CamelCase),
	resp.WithDefaultLocale("en"),
	resp.WithErrorFormat(resp.ProblemDetails),
)
if err != nil {
	return err
}

return api.OK(ctx, data)
```

Engine yang dibuat bersifat immutable setelah pembuatan dan aman dipakai bersama. Pilihan naming yang tersedia: `resp.SnakeCase`, `resp.CamelCase`, dan `resp.PascalCase`. Format error: `resp.Envelope` atau `resp.ProblemDetails`.

## Konfigurasi default

`resp.Configure` mengganti engine default secara atomik:

```go
if err := resp.Configure(resp.WithDefaultLocale("en")); err != nil {
	log.Fatal(err)
}
```

Panggil sekali saat aplikasi mulai. Untuk pengujian, gunakan engine khusus agar konfigurasi global tidak saling memengaruhi.

## Opsi response

`resp.Locale("en")` dan `resp.Location("/users/123")` hanya berlaku untuk pemanggilan response tersebut. `Location` diterapkan pada response `201 Created`.

Opsi internal error sengaja hanya tersedia di tingkat engine; opsi itu bukan pilihan per request atau header.
