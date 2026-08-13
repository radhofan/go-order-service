# Go Order Service

Service backend REST API sederhana menggunakan bahasa pemrograman Go dan SQLite untuk mengelola produk, customer, serta pemrosesan pesanan.

## Prasyarat

- Go versi 1.21 atau lebih baru

## Struktur Proyek

```
GoOrderService/
├── api-demo/          # Berkas pengujian request HTTP sample
├── backend/
│   ├── cmd/api/       # Entry point utama aplikasi
│   ├── internal/
│   │   ├── config/    # Konfigurasi environment
│   │   ├── db/        # Inisialisasi database dan eksekusi schema.sql
│   │   ├── domain/    # Struct model, request, response, dan error kustom
│   │   ├── handler/   # HTTP handlers dan penanganan response JSON
│   │   ├── middleware/# Logging middleware untuk setiap request HTTP
│   │   ├── repository/# Access layer basis data SQLite
│   │   └── service/   # Aturan bisnis utama dan transaksi atomik
│   ├── schema.sql     # Skema tabel database SQLite
│   └── test/          # Skenario pengujian unit test dan integrasi
└── Studi_Kasus_Junior_Developer-1-4.pdf
```

## Cara Menjalankan Aplikasi

1. Masuk ke direktori backend

   ```bash
   cd backend
   ```

2. Jalankan server HTTP
   ```bash
   go run ./cmd/api
   ```
   Server akan berjalan pada port `8080` (atau sesuai konfigurasi variabel `PORT` pada `.env`).

## Cara Menjalankan Pengujian (Unit Test)

Jalankan seluruh suite unit test

```bash
cd backend
go test -v ./test
```

## Daftar Lengkap API Endpoints

### 1. Produk (Products)

- **POST /products**
  Membuat produk baru.
  Request Body: `{"sku": "KAOS-01", "name": "Kaos Hitam", "price": 85000, "stock": 20}`
  Response: `201 Created`

- **GET /products**
  Mengambil daftar produk dengan fitur pencarian dan paginasi.
  Query Parameters: `?page=1&limit=10&q=kaos`
  Response: `200 OK` (disertai metadata total data)

- **GET /products/{id}**
  Mengambil detail produk berdasarkan ID.
  Response: `200 OK` atau `404 Not Found`

- **PUT /products/{id}**
  Mengubah data produk (name, price, stock).
  Request Body: `{"name": "Kaos Hitam Premium", "price": 95000, "stock": 15}`
  Response: `200 OK` atau `404 Not Found`

### 2. Customer

- **POST /customers**
  Membuat customer baru dengan validasi format email dan keunikan email.
  Request Body: `{"name": "Budi Santoso", "email": "budi@example.com"}`
  Response: `201 Created` atau `409 Conflict` (jika email sudah terdaftar)

### 3. Pesanan (Orders)

- **POST /orders**
  Membuat pesanan baru. Total harga dihitung otomatis oleh server dari harga produk saat pesanan dibuat.
  Request Body: `{"customer_id": 1, "items": [{"product_id": 1, "qty": 2}]}`
  Response: `201 Created`, `400 Bad Request` (stok kurang), `404 Not Found` (customer/produk tidak ditemukan)

- **GET /orders/{id}**
  Mengambil detail pesanan lengkap beserta rincian item, nama produk, harga saat dipesan (`price_at_order`), dan total harga.
  Response: `200 OK` atau `404 Not Found`

- **GET /orders**
  Mengambil daftar seluruh pesanan. Bisa difilter berdasarkan customer atau status.
  Query Parameters: `?customer_id=1&status=PENDING`
  Response: `200 OK`

- **PATCH /orders/{id}/status**
  Mengubah status pesanan mengikuti alur validasi.
  Request Body: `{"status": "PAID"}`
  Response: `200 OK` atau `409 Conflict` (jika alur transisi status tidak valid)

## Notes

- **Penanganan Produk/Customer Tidak Ditemukan (404 Not Found)**: Saat pembuatan pesanan (`POST /orders`), jika `customer_id` atau `product_id` tidak ada di database, sistem mengembalikan `404 Not Found` (bukan 400). Alasannya adalah semantik HTTP standar RFC 7231 di mana entitas identitas yang dirujuk tidak ditemukan di database.
- **Perhitungan Total Server-Side**: Client tidak mengirimkan harga. Harga diambil langsung dari tabel `products` saat pesanan dibuat dan disimpan pada `order_items.price_at_order` agar perubahan harga produk di kemudian hari tidak mengubah data historis transaksi.
- **Transaksi Atomik Database**: Pembuatan pesanan dan pengurangan stok produk dijalankan dalam satu transaksi database (`BEGIN...COMMIT`). Jika salah satu item gagal (stok kurang atau produk tidak ada), seluruh transaksi dibatalkan (`ROLLBACK`).
- **Alur Transisi Status Pesanan**: Alur status dibatasi ketat: `PENDING -> PAID -> SHIPPED -> COMPLETED` atau `CANCELLED`. Transisi yang melompati atau mundur dari urutan ini dikembalikan sebagai `409 Conflict`.
- **Pengembalian Stok Pada Pembatalan**: Mengubah status pesanan menjadi `CANCELLED` secara otomatis dan atomik mengembalikan seluruh stok item pesanan kembali ke tabel `products`.
