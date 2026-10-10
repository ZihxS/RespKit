# Response HTTP

## Helper dasar

| Helper | Status | Kegunaan |
| --- | ---: | --- |
| `resp.OK(ctx, data)` | 200 | Response berhasil |
| `resp.Created(ctx, data, opts...)` | 201 | Resource baru berhasil dibuat |
| `resp.NoContent(ctx)` | 204 | Sukses tanpa body |
| `resp.Fail(ctx, err, opts...)` | Mengikuti klasifikasi error | Response gagal |
| `resp.NotFound(ctx, opts...)` | 404 | Shortcut untuk resource tidak ditemukan |

Semua helper mengembalikan `error`. Tangani error itu sesuai framework. Error tersebut menandakan konteks tidak didukung atau proses penulisan response gagal; jangan mencoba menulis response kedua jika response pertama mungkin sudah terkirim.

## Bentuk response default

Response sukses:

```json
{
  "success": true,
  "data": {"user_id": 123}
}
```

Response gagal menggunakan kode dan pesan yang aman untuk client:

```json
{
  "success": false,
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "Data tidak ditemukan."
  }
}
```

Kunci milik RespKit memakai `snake_case` secara default. Pilihan casing hanya mengubah kunci envelope RespKit; JSON `data` milik aplikasi tidak ditulis ulang.

## Header dan perilaku HTTP

- `Content-Type` memakai JSON UTF-8. Format Problem Details menggunakan `application/problem+json`.
- `Content-Language` menyatakan bahasa pesan response; `Vary: Accept-Language` membantu cache yang memperhatikan locale.
- Header `X-Request-ID` yang aman dapat diteruskan ke metadata response.
- `resp.Created` mendukung `resp.Location(url)` untuk header `Location`.
- Response `204 No Content` dan request `HEAD` tidak mengirim body.
- Response `401` menyertakan challenge `WWW-Authenticate: Bearer` secara default. Aplikasi dapat mengatur challenge yang sesuai skema autentikasinya.

Contoh:

```go
err := resp.Created(ctx, user, resp.Location("/users/123"))
```

## RFC 9457 Problem Details

Format error dapat diubah per engine:

```go
api, err := resp.New(resp.WithErrorFormat(resp.ProblemDetails))
if err != nil {
	return err
}

return api.Fail(ctx, resp.ErrNotFound)
```

Success response tetap menggunakan envelope RespKit; opsi tersebut mengubah representasi error.
