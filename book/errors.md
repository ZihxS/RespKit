# Pemetaan error

RespKit memetakan error berdasarkan tipe dan identitasnya, termasuk error yang dibungkus dengan `%w`. Pencocokan menggunakan `errors.Is` dan `errors.As`, bukan membaca teks `Error()`.

## Kode bawaan

| Error domain | Kode publik | HTTP |
| --- | --- | ---: |
| `resp.ErrBadRequest` | `BAD_REQUEST` | 400 |
| `resp.ErrUnauthorized` | `UNAUTHORIZED` | 401 |
| `resp.ErrForbidden` | `FORBIDDEN` | 403 |
| `resp.ErrNotFound` | `RESOURCE_NOT_FOUND` | 404 |
| `resp.ErrConflict` | `RESOURCE_CONFLICT` | 409 |
| `resp.ErrValidation` | `VALIDATION_FAILED` | 422 |
| `resp.ErrRateLimited` | `RATE_LIMITED` | 429 |
| `resp.ErrUnavailable` | `SERVICE_UNAVAILABLE` | 503 |
| Error tak dikenal | `INTERNAL_ERROR` | 500 |

Error sintaks atau tipe JSON yang tidak sesuai dipetakan ke `BAD_REQUEST`. SQLSTATE `23505` dikenali sebagai pelanggaran unik dan menjadi `RESOURCE_ALREADY_EXISTS` (409). Untuk kondisi database lain, tentukan makna domain secara eksplisit atau gunakan [adapter pgx](postgres.md).

## Membuat error aplikasi

```go
err := resp.NewError(
	"PAYMENT_DECLINED",
	http.StatusUnprocessableEntity,
	resp.WithCause(errFromPaymentProvider),
)
return resp.Fail(ctx, err)
```

Kode publik harus stabil, berupa huruf kapital, angka, dan underscore. `Cause` membantu logging dan inspeksi di server, tetapi tidak diserialisasi ke response. Jika kode/status tidak valid, resolver memakai fallback internal yang aman.

## Detail validasi

```go
err := resp.NewError(
	"VALIDATION_FAILED",
	422,
	resp.WithDetails(
		resp.FieldError{Field: "email", Code: "FIELD_REQUIRED"},
		resp.FieldError{
			Field:  "name",
			Code:   "MIN_LENGTH",
			Params: map[string]string{"min": "3"},
		},
	),
)
return resp.Fail(ctx, err)
```

Detail membantu client menandai field yang perlu diperbaiki. Hanya masukkan parameter yang aman dan memang ditujukan untuk client; jangan taruh token, query, atau data sensitif di sana.

## Mapper custom

Gunakan mapper saat tipe error pihak ketiga mempunyai arti domain yang sudah dipahami aplikasi. Mapper dapat mengembalikan `false` jika tidak cocok.

```go
func mapPaymentError(err error) (resp.Definition, bool) {
	var paymentErr *PaymentError
	if !errors.As(err, &paymentErr) {
		return resp.Definition{}, false
	}

	return resp.Definition{
		Code:   "PAYMENT_DECLINED",
		Status: http.StatusUnprocessableEntity,
	}, true
}

api, err := resp.New(resp.WithMapper("payment", mapPaymentError))
if err != nil {
	return err
}
return api.Fail(ctx, errFromPaymentProvider)
```

Mapper dikonfigurasi saat pembuatan engine. Jangan memetakan error yang ambigu hanya berdasarkan teks.
