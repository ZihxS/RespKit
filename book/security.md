# Keamanan

Raw `err.Error()` tidak pernah dimasukkan ke pesan publik secara default. Error yang tidak dikenal mendapat kode umum `INTERNAL_ERROR` dan status 500. Gunakan kode error domain yang stabil, dan jangan memakai pesan mentah database/provider sebagai teks untuk client.

## Jangan aktifkan debug message di production

`resp.WithExposeInternalErrors(true)` menambahkan `debug_message` berisi pesan error asli. Pesan itu bisa mengandung query SQL, token, jalur internal, alamat host, atau informasi pribadi.

```go
// Hanya untuk proses lokal yang terisolasi.
api, err := resp.New(resp.WithExposeInternalErrors(true))
```

Jangan aktifkan opsi ini pada server yang dapat diakses melalui internet. Jangan mengaktifkannya berdasarkan header request atau input pengguna.

## Pelaporan celah keamanan

Laporkan dugaan kerentanan secara privat ke [m.saleh.solahudin@gmail.com](mailto:m.saleh.solahudin@gmail.com). Jangan membuat issue publik berisi kredensial, data pengguna, log production, atau langkah eksploitasi.
