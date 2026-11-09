# jsonapi - Manual Penggunaan

Bagian dari ekosistem [GoCroot](https://gocroot.if.co.id/). Dipakai oleh [GoCroot Core](https://github.com/gocroot/core).

## Pendahuluan
`jsonapi` adalah paket Go yang menyediakan fungsi untuk melakukan HTTP request dengan JSON payload dan autentikasi menggunakan token. Paket ini mendukung metode GET dan POST dengan opsi untuk menambahkan token otorisasi.

## Instalasi
Tambahkan paket ini ke proyek Go Anda:

```sh
 go get github.com/gocroot/jsonapi
```

## Fungsi-Fungsi

### 1. `PostStructWithToken`
Mengirim data dalam bentuk struct ke endpoint dengan metode POST dan menambahkan token otorisasi ke dalam header.

#### Parameter:
- `tokenkey` (string): Nama header token.
- `tokenvalue` (string): Nilai token.
- `structname` (interface{}): Data yang akan dikirim dalam format JSON.
- `urltarget` (string): URL tujuan.

#### Contoh Penggunaan:
```go
var response map[string]interface{}
data := map[string]string{"key": "value"}
status, result, err := jsonapi.PostStructWithToken("Authorization", "Bearer mytoken", data, "https://example.com/api")
```

### 2. `Get`
Melakukan request GET ke endpoint yang diberikan dan mengembalikan hasil JSON.

#### Parameter:
- `urltarget` (string): URL tujuan.

#### Contoh Penggunaan:
```go
var response map[string]interface{}
status, result, err := jsonapi.Get("https://example.com/api")
```

### 3. `GetWithBearer`
Melakukan request GET dengan menambahkan token bearer ke dalam header.

#### Parameter:
- `tokenbearer` (string): Token otorisasi.
- `urltarget` (string): URL tujuan.

#### Contoh Penggunaan:
```go
var response map[string]interface{}
status, result, err := jsonapi.GetWithBearer("mytoken", "https://example.com/api")
```

### 4. `GetStructWithToken`
Melakukan request GET dengan menambahkan token kustom ke dalam header.

#### Parameter:
- `tokenkey` (string): Nama header token.
- `tokenvalue` (string): Nilai token.
- `urltarget` (string): URL tujuan.

#### Contoh Penggunaan:
```go
var response map[string]interface{}
status, result, err := jsonapi.GetStructWithToken("X-API-KEY", "myapikey", "https://example.com/api")
```

## Error Handling
Jika response bukan JSON yang valid, fungsi akan mengembalikan error dengan potongan konten yang diterima (maksimal 512 byte). Password di URL disamarkan.

## Keamanan & Performa
- Semua fungsi memakai satu `jsonapi.Client` bersama (keep-alive) dengan timeout 30 detik. Bisa diganti: `jsonapi.Client = &http.Client{Timeout: 10 * time.Second}`.
- Body respons dibatasi `jsonapi.MaxResponseBytes` (default 10 MiB). Respons yang lebih besar ditolak dengan error.
- Error dari `json.Marshal` pada `PostStructWithToken` dikembalikan, tidak lagi mengirim body kosong.

## Lisensi
Paket ini tersedia dengan lisensi open-source. Silakan gunakan sesuai kebutuhan Anda.

## Release Version

Rilis dibuat **otomatis** oleh workflow *Auto Tag & Release* setiap ada push ke `main` yang mengubah kode Go, `go.mod`, atau `go.sum`. Jangan membuat tag manual. Kenaikan versi mengikuti pesan commit: `fix:` → patch, `feat:` → minor, dan `[skip release]` untuk melewati rilis. Rilis manual bisa lewat tab Actions → *Run workflow*.

```sh
go get -u ./...
go mod tidy
go list -m -versions github.com/gocroot/jsonapi
```
