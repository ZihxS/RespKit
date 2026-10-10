# CLI diagnostics

Perintah di `cmd/respkit` memberikan informasi tentang versi, katalog error, dan klasifikasi tanpa memerlukan API key atau mengirim data ke layanan AI.

```bash
go run ./cmd/respkit version
go run ./cmd/respkit doctor
go run ./cmd/respkit catalog
go run ./cmd/respkit explain RESOURCE_NOT_FOUND
```

Diagnostics error tersedia di Go melalui `resp.Inspect(err)` dan `resp.Explain(err)`. Hasilnya berisi metadata klasifikasi seperti mapper, kode sumber, kode publik, dan status HTTP. Hindari mencatat atau membagikan error mentah yang mungkin memuat data sensitif.
