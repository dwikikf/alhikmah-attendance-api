## Relevant Files

- `internal/handler/class_handler.go` - Entry point API untuk data kelas, butuh meneruskan `teacherID` jika login sebagai guru.
- `internal/repository/class_postgres.go` - Perlu memfilter *query* SQL kelas berdasarkan `teacher_id`.
- `internal/handler/student_handler.go` - Entry point API untuk siswa, butuh validasi `teacherID`.
- `internal/repository/student_postgres.go` - Perlu melakukan *JOIN* ke tabel `classes` secara dinamis hanya jika memfilter berdasarkan `teacherID`.

### Notes
- Pastikan penggunaan *JOIN* di `student_postgres.go` dilakukan dinamis (hanya jika `teacher_id` tidak kosong) agar tidak memperberat eksekusi untuk Admin.
- Validasi CRUD akan menggunakan pengecekan tambahan apakah kelas/siswa yang diedit berada di bawah kepemilikan guru tersebut.

## Instructions for Completing Tasks

**IMPORTANT:** As you complete each task, you must check it off in this markdown file by changing `- [ ]` to `- [x]`. This helps track progress and ensures you don't skip any steps.

## Tasks

- [x] 0.0 Create feature branch
  - [x] 0.1 Checkout branch baru `feature/teacher-data-access-v3` di repositori API.
  - [x] 0.2 Checkout branch baru `feature/teacher-data-access-v3` di repositori Web.
- [x] 1.0 Terapkan Filter & Validasi Kelas (Class API) untuk Role Teacher
  - [x] 1.1 Modifikasi `class_handler.go` (`GetAll`) untuk memeriksa `role == "teacher"` (menggantikan "guru") dan meneruskan `teacherID`.
  - [x] 1.2 Pastikan `class_postgres.go` (`GetAll`) mengeksekusi parameter `teacherID` yang diteruskan.
  - [x] 1.3 Tambahkan validasi kepemilikan `CheckClassOwnership` untuk method `Update` dan `Delete` di handler kelas.
- [x] 2.0 Terapkan Filter & Validasi Siswa (Student API) untuk Role Teacher secara Dinamis
  - [x] 2.1 Modifikasi `student_handler.go` (`GetAll`) untuk menangkap role `"teacher"` dan mengirim `teacherID` ke service.
  - [x] 2.2 Modifikasi `student_postgres.go` (`GetAll`) untuk menyisipkan `LEFT JOIN` pada `countQuery` HANYA jika `teacherID != ""`, demi performa yang sangat ringan.
  - [x] 2.3 Inject `ClassService` ke `StudentHandler` dan periksa kepemilikan kelas saat `teacher` memicu aksi `Create`, `Update`, atau `Delete` pada siswa.
- [x] 3.0 Verifikasi Keseluruhan API
  - [x] 3.1 Kompilasi ulang *backend* dengan `go build -v ./cmd/api`.
  - [x] 3.2 Pastikan respons data siswa & kelas di Web menjadi benar tanpa menyentuh *source code* Web (karena data sudah terfilter murni di *backend*).
