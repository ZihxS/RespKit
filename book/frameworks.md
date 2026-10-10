# HTTP framework

RespKit mempunyai satu package API:

```go
import "github.com/ZihxS/RespKit/resp"
```

Context Gin, Echo v4/v5, dan Fiber v2/v3 dapat langsung diberikan ke helper RespKit. Untuk `net/http`, Chi, Gorilla/mux, dan Beego v2, bungkus writer dan request dengan `resp.HTTP(w, r)`.

## net/http, Chi, Gorilla/mux, dan Beego

```go
func handler(w http.ResponseWriter, r *http.Request) {
	if err := resp.OK(resp.HTTP(w, r), map[string]any{"user_id": 7}); err != nil {
		log.Printf("write response: %v", err)
	}
}
```

Fungsi handler tersebut bisa didaftarkan pada router Chi atau Gorilla/mux karena keduanya memakai `net/http`. Untuk Beego v2, gunakan `resp.HTTP(c.ResponseWriter, c.Request)` dari context Beego.

## Gin

```go
func handler(c *gin.Context) {
	if err := resp.OK(c, user); err != nil {
		_ = c.Error(err)
	}
}
```

## Echo v4 dan v5

Echo mengembalikan error dari handler ke error handler framework:

```go
func handler(c echo.Context) error { // Echo v4
	return resp.Fail(c, resp.ErrNotFound)
}
```

Pada Echo v5, tipe context handler adalah `*echo.Context`; gunakan package import v5.

## Fiber v2 dan v3

Kembalikan error yang dihasilkan helper kepada Fiber:

```go
func handler(c *fiber.Ctx) error { // Fiber v2
	return resp.OK(c, user)
}
```

Fiber v3 memakai tipe `fiber.Ctx` dari package v3.

## Batas kompatibilitas

Framework tidak menjadi dependency modul utama. Kesesuaian tiap versi diuji pada router upstream secara terpisah; versi tersebut bukan janji bahwa semua versi framework lama/baru sudah diverifikasi. Jalankan integration suite pada versi yang dipakai aplikasi.
