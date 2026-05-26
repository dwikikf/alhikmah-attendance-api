## Relevant Files

- `internal/handler/report_handler.go` - Implementasi fungsi `Export`.
- `internal/service/report_service.go` - Logika penarikan data mentah untuk kemudian diubah menjadi format ekspor.
- (Frontend) `src/components/reports/ReportExporter.tsx` - Hapus opsi menu PDF.

### Notes
- *Endpoint* `POST /reports/export` sudah diregistrasi di `main.go`.
- Frontend sudah memiliki komponen `ReportExporter` yang mengirimkan *payload* berisikan jenis laporan (*daily*, *monthly*, *semester*) dan format yang diinginkan (*excel*, *csv*).
- Backend akan mem-_parsing_ payload tersebut dan merangkai fail CSV atau XLSX.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Checkout branch baru `feature/report-export` di repositori API.
  - [x] 0.2 Checkout branch baru `feature/report-export` di repositori Web.
- [x] 1.0 Sesuaikan UI Frontend (Web)
  - [x] 1.1 Hapus elemen `<DropdownMenuItem>` untuk PDF dari komponen `ReportExporter.tsx`.
- [x] 2.0 Implementasi Fungsi Ekspor Excel & CSV (Backend API)
  - [x] 2.1 Buka `internal/handler/report_handler.go`.
  - [x] 2.2 Buat struktur `ExportRequest` untuk menerima format JSON payload.
  - [x] 2.3 Rangkai struktur kerangka CSV (menggunakan `encoding/csv`) dengan mengambil data report dari `ReportService`.
  - [x] 2.4 Rangkai struktur kerangka Excel (menggunakan `github.com/xuri/excelize/v2`) dengan mengambil data report dari `ReportService`.
  - [x] 2.5 Gabungkan ke dalam kontrol *switch-case* (CSV / Excel) beserta validasi *header response*.
- [x] 3.0 Pengujian API dengan Frontend
  - [x] 3.1 Kompilasi ulang backend API (`go build`).
  - [x] 3.2 Lakukan klik tombol unduh dari Web (CSV & Excel) dan pastikan berkas berisikan tabel data absensi utuh.
