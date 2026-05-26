## Relevant Files

- `internal/handler/class_handler.go` - Akan ditambahkan fungsi untuk *Export Excel* dan *Export QR Code ZIP*.
- `internal/service/class_service.go` - Logika *business* untuk merangkai data siswa dan mengubahnya menjadi Excel & ZIP.
- `cmd/api/main.go` - Menambahkan *routing* API baru (`/classes/:class_id/export-excel` dan `/classes/:class_id/export-qrcode`).
- `src/components/classes/ClassDetail.tsx` (Frontend Web) - Tombol *Export Excel* dan *Download QR Code (ZIP)* akan ditambahkan di sini.

### Notes
- Penggunaan library `archive/zip` (bawaan Go) dan library QR Code (misalnya `skip2/go-qrcode`) untuk mem-*generate* file `.png` lalu memasukkannya ke dalam `.zip`.
- Untuk Excel, bisa menggunakan package populer seperti `github.com/xuri/excelize`.
- Format unduhan QR akan berupa kumpulan gambar (seperti `NISN_NamaSiswa.png`) dalam 1 arsip `.zip`.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Checkout branch baru `feature/class-export-tools` di repositori API.
  - [x] 0.2 Checkout branch baru `feature/class-export-tools` di repositori Web.
- [x] 1.0 Buat Backend API Endpoint untuk Export Excel Data Kelas
  - [x] 1.1 `go get github.com/xuri/excelize/v2` di repositori API.
  - [x] 1.2 Buat method `ExportExcel` pada `ClassHandler` yang mengambil daftar siswa melalui *StudentService* atau *Repository*.
  - [x] 1.3 Format data siswa ke dalam kerangka dokumen `.xlsx` baru.
  - [x] 1.4 Daftarkan *routing* `GET /api/v1/classes/:class_id/export-excel` di `cmd/api/main.go` (termasuk *role middleware* guru/admin).
- [x] 2.0 Buat Backend API Endpoint untuk Export QR Code (Format ZIP Berisi Gambar)
  - [x] 2.1 `go get github.com/skip2/go-qrcode` di repositori API.
  - [x] 2.2 Buat method `ExportQRCode` pada `ClassHandler`.
  - [x] 2.3 Rangkai logika iterasi siswa, konversi `QRCodeData` masing-masing siswa menjadi *stream buffer* gambar PNG (ukuran ~256x256).
  - [x] 2.4 Kemas gambar-gambar tersebut ke dalam berkas arsip `.zip` dan kembalikan *header response* yang tepat (`application/zip`).
  - [x] 2.5 Daftarkan *routing* `GET /api/v1/classes/:class_id/export-qrcode` di `cmd/api/main.go`.
- [x] 3.0 Sesuaikan Frontend Web untuk Tombol Unduh di Halaman Detail Kelas
  - [x] 3.1 Buka `src/components/classes/ClassDetail.tsx` di *Web workspace*.
  - [x] 3.2 Tambahkan 2 tombol (*Export Excel* & *Download QR (ZIP)*) pada *Header* (*CardTitle*) atau *Toolbar*.
  - [x] 3.3 Buat fungsi `handleExportExcel` & `handleExportQR` yang memanfaatkan utilitas `window.URL.createObjectURL` untuk memicu aksi unduh lokal dari hasil balasan *Blob API*.
  - [x] 3.4 Validasi kompilasi akhir Web dan uji fitur ekspor secara keseluruhan.
