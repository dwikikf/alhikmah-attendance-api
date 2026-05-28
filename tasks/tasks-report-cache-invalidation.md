## Relevant Files

- `internal/domain/report.go` - Interface domain `ReportService` dan `ReportCacheRepository` yang perlu diupdate dengan menambahkan method `ForceRefresh`.
- `internal/service/report_service.go` - Logika bisnis utama report. Perlu diubah agar menerima flag `forceRefresh` dan menghapus cache lama sebelum regenerasi.
- `internal/repository/report_cache_postgres.go` - Repository cache tabel `reports`. Perlu ditambahkan method `Delete` untuk menghapus entri cache lama saat force refresh.
- `internal/handler/report_handler.go` - Handler HTTP. Perlu membaca query param `?force_refresh=true` dan meneruskannya ke service.
- `internal/dto/` - (Opsional) Jika ingin menambahkan DTO baru untuk response yang menyertakan metadata cache (`is_cached`, `cached_at`).

### Notes

- Tabel `reports` sudah ada di database dan sudah berfungsi sebagai cache dengan TTL 1 jam (lihat query `WHERE generated_at > NOW() - INTERVAL '1 hour'` di `report_cache_postgres.go`).
- Cache saat ini **sudah berjalan** untuk semua jenis laporan (harian, bulanan, semesteran). Yang belum ada adalah mekanisme **force refresh** dari luar (via API).
- Perubahan ini bersifat **backward-compatible** — `force_refresh` adalah query param opsional, jadi semua client yang tidak mengirimkan parameter ini akan tetap berperilaku seperti sebelumnya.
- Tidak ada unit test yang perlu ditulis karena proyek ini tidak menggunakan framework testing.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

Example:
- `- [ ] 1.1 Read file` → `- [x] 1.1 Read file` (after completing)

Update the file after completing each sub-task, not just after completing an entire parent task.

## Tasks

- [ ] 0.0 Create feature branch
  - [ ] 0.1 Buat dan checkout branch baru: `git checkout -b feature/report-cache-invalidation`

- [ ] 1.0 Tambahkan method `Delete` pada `ReportCacheRepository`
  - [ ] 1.1 Buka file `internal/domain/report.go` dan tambahkan method baru `Delete(reportType, classID, periodStart, periodEnd string) error` ke dalam interface `ReportCacheRepository`
  - [ ] 1.2 Buka file `internal/repository/report_cache_postgres.go` dan implementasikan method `Delete` dengan query `DELETE FROM reports WHERE report_type = $1 AND class_id = $2 AND period_start = $3 AND period_end = $4`

- [ ] 2.0 Ubah `ReportService` agar mendukung flag `forceRefresh`
  - [ ] 2.1 Buka file `internal/domain/report.go` dan ubah signature semua method dalam interface `ReportService` agar menerima parameter tambahan `forceRefresh bool`:
    - `GetDailyReport(classID, dateStr string, forceRefresh bool) (*DailyReport, error)`
    - `GetMonthlyReport(classID, monthStr string, forceRefresh bool) (*MonthlyReport, error)`
    - `GetSemesterReport(classID, academicYear string, semester int, forceRefresh bool) (*SemesterReport, error)`
    - `GetStudentReport(studentID, startDate, endDate string) (*StudentReport, error)` — ini tidak perlu diubah karena laporan per-siswa tidak di-cache
  - [ ] 2.2 Buka file `internal/service/report_service.go` dan update fungsi `GetDailyReport`:
    - Tambahkan parameter `forceRefresh bool`
    - Jika `forceRefresh == true`, panggil `s.cache.Delete(...)` sebelum melakukan pengecekan cache
    - Logika pengambilan data dan penyimpanan cache tetap sama seperti sekarang
  - [ ] 2.3 Update fungsi `GetMonthlyReport` di service dengan pola yang sama seperti 2.2
  - [ ] 2.4 Update fungsi `GetSemesterReport` di service dengan pola yang sama seperti 2.2

- [ ] 3.0 Update handler untuk membaca query param `force_refresh`
  - [ ] 3.1 Buka file `internal/handler/report_handler.go` dan update `GetDailyReport`:
    - Baca query param: `forceRefresh := c.Query("force_refresh") == "true"`
    - Teruskan ke service: `h.service.GetDailyReport(classID, dateStr, forceRefresh)`
  - [ ] 3.2 Update `GetMonthlyReport` di handler dengan pola yang sama seperti 3.1
  - [ ] 3.3 Update `GetSemesterReport` di handler dengan pola yang sama seperti 3.1
  - [ ] 3.4 Update fungsi `Export` di handler — karena `Export` juga memanggil `service.GetDailyReport`, `GetMonthlyReport`, dan `GetSemesterReport` secara internal, tambahkan `forceRefresh: false` (hardcode false, karena export selalu mengambil data terbaru dari cache yang sudah ada, tidak perlu force refresh dari sini)

- [ ] 4.0 Update Frontend (`alhikmah-attendance-web`) untuk mendukung force refresh
  - [ ] 4.1 Temukan file `useQuery` atau fungsi fetcher yang mengambil data laporan bulanan dan semesteran di workspace `alhikmah-attendance-web`
  - [ ] 4.2 Tambahkan state `isRefreshing` (boolean) di komponen halaman laporan
  - [ ] 4.3 Buat fungsi `handleForceRefresh` yang memanggil API dengan tambahan query param `?force_refresh=true` menggunakan `useMutation` atau `refetch` dari React Query
  - [ ] 4.4 Tambahkan tombol "🔄 Perbarui Data" di UI halaman laporan (misalnya di sebelah dropdown bulan/semester) yang memanggil `handleForceRefresh` saat ditekan
  - [ ] 4.5 Tambahkan tampilan loading state saat proses force refresh berlangsung dan tampilkan notifikasi sukses/gagal menggunakan `toast`
  - [ ] 4.6 Setelah `handleForceRefresh` berhasil, panggil `queryClient.invalidateQueries(...)` agar `useQuery` utama melakukan refetch dengan data yang sudah diperbarui

- [ ] 5.0 Verifikasi dan testing manual
  - [ ] 5.1 Jalankan backend (`go run ./cmd/...` atau sesuai perintah proyek) dan pastikan tidak ada error kompilasi
  - [ ] 5.2 Test endpoint daily report tanpa force refresh: `GET /reports?class_id=...&date=...` — pastikan cache lama masih terpakai (respons cepat)
  - [ ] 5.3 Test endpoint daily report dengan force refresh: `GET /reports?class_id=...&date=...&force_refresh=true` — pastikan data diambil ulang dari DB dan cache diperbarui (cek kolom `generated_at` di tabel `reports` berubah)
  - [ ] 5.4 Test endpoint monthly report dan semester report dengan skenario yang sama
  - [ ] 5.5 Test di frontend: tekan tombol "Perbarui Data" dan pastikan data di tabel/grafik ikut berubah jika ada perubahan absensi yang baru di-input

- [ ] 6.0 Commit dan push
  - [ ] 6.1 Commit semua perubahan backend: `git commit -m "feat: add force_refresh support for report cache invalidation"`
  - [ ] 6.2 Commit semua perubahan frontend di workspace `alhikmah-attendance-web`: `git commit -m "feat: add refresh button for report cache invalidation"`
  - [ ] 6.3 Push kedua branch ke remote
