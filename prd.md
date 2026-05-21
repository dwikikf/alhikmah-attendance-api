# PRD: Aplikasi Absensi Siswa SD dengan QR Code

## Overview

Aplikasi absensi siswa SD berbasis QR code yang memudahkan pencatatan kehadiran siswa secara real-time. Setiap siswa memiliki QR code unik yang digenerate dari NISN (Nomor Induk Siswa Nasional). QR code di-scan menggunakan kamera device oleh guru pada saat awal pembelajaran di kelas. Sistem mencatat data kehadiran siswa dan menyediakan laporan absensi harian, mingguan, dan bulanan bahkan per-semester paling penting untuk guru.

**Masalah yang diselesaikan:**

- Proses absensi manual memakan waktu dan rentan kesalahan pencatatan
- Kesulitan dalam verifikasi identitas siswa
- Tidak ada track record digital yang mudah diakses dan dilaporkan

---

## Goals

1. Mempercepat proses pencatatan kehadiran siswa dari manual (~5-10 menit) menjadi digital (~1-2 menit)
2. Mengurangi human error dalam pencatatan absensi
3. Menyediakan data kehadiran digital yang dapat diakses dan dilaporkan dengan mudah
4. Memudahkan guru untuk memantau kehadiran siswa real-time di dalam kelas

---

## User Stories

### User Story 1: Guru Mencatat Absensi

**Sebagai seorang guru**, saya ingin dapat scan QR code siswa dengan cepat pada awal pembelajaran **sehingga** pencatatan kehadiran menjadi efisien dan akurat.

**Acceptance Criteria:**

- Guru membuka aplikasi dan memilih kelas + tanggal
- Guru siap dengan device (smartphone/tablet)
- Setiap siswa menunjukkan QR code untuk di-scan
- QR code berhasil di-scan dan data terekam

### User Story 2: Guru Melihat Laporan Absensi Harian

**Sebagai seorang guru**, saya ingin melihat laporan absensi harian siswa **sehingga** saya dapat mengetahui siswa mana yang hadir dan mana yang absen.

**Acceptance Criteria:**

- Guru dapat membuka laporan per tanggal
- Laporan menampilkan daftar siswa (hadir/absen)
- Laporan dapat diunduh atau dibagikan

### User Story 3: Siswa Scan QR Code Sendiri

**Sebagai seorang siswa**, saya ingin scan QR code saya sendiri dengan mudah **sehingga** saya dapat mencatat kehadiran saya sendiri.

**Acceptance Criteria:**

- Interface simpel dan intuitif untuk siswa
- QR code saya dapat di-scan dengan kamera device
- Konfirmasi kehadiran ditampilkan setelah scan berhasil

### User Story 4: Guru Input Absensi Manual

**Sebagai seorang guru**, saya ingin dapat input data absensi manual untuk siswa yang tidak hadir dengan alasan izin atau sakit **sehingga** data absensi tetap lengkap dan akurat meskipun siswa tidak bisa scan.

**Acceptance Criteria:**

- Guru dapat membuka daftar siswa kelas yang belum ter-scan
- Guru dapat memilih satu atau multiple siswa sekaligus
- Guru dapat menentukan status kehadiran (Hadir, Izin, Sakit, Tanpa Keterangan)
- Guru dapat menambahkan catatan/keterangan (opsional)
- Data absensi manual terekam dengan timestamp dan nama guru yang input

---

## Functional Requirements

### FR1: QR Code Generation

1. Sistem harus dapat generate QR code unik untuk setiap siswa berdasarkan NISN
2. QR code harus berisi data siswa minimum: NISN, Nama Siswa, Kelas
3. QR code dapat di-export sebagai gambar (PNG/JPG) untuk dicetak atau ditampilkan digital
4. Setiap QR code harus dapat di-scan dan divalidasi

### FR2: Scanning & Validation

1. Aplikasi harus membuka akses ke kamera device (smartphone/tablet)
2. Aplikasi harus dapat scan QR code dan ekstrak data siswa
3. Sistem harus validasi NISN dan identitas siswa sebelum mencatat kehadiran
4. Jika QR code tidak valid atau tidak dikenali, sistem harus menampilkan pesan error
5. Setiap scan harus mencatat timestamp kehadiran

### FR3: Attendance Recording (QR Code Scanning)

1. Sistem harus mencatat waktu scan QR code sebagai bukti kehadiran
2. Setiap siswa hanya bisa di-scan 1x per hari per kelas (prevent duplikat)
3. Guru dapat melihat daftar siswa yang sudah ter-scan (real-time)
4. Tanggal dan kelas secara otomatis terambil dari data siswa di QR code, tidak perlu input manual
5. Sistem mencegah duplicate scan dengan notifikasi "Siswa sudah ter-scan" jika QR code di-scan ulang

### FR3B: Manual Attendance Entry

1. Guru dapat membuka fitur "Input Manual" untuk entry data kehadiran siswa
2. Guru dapat memilih satu atau multiple siswa dari daftar kelas
3. Guru dapat set status: Hadir, Izin, Sakit, Tanpa Keterangan
4. Setiap entry manual harus mencatat: nama guru, waktu input, keterangan (opsional)
5. Sistem dapat mengubah status jika sudah ter-input (dengan audit trail)
6. Guru tidak bisa menghapus, hanya mengubah status absensi siswa

### FR4: Data Management

1. Aplikasi harus menyimpan data absensi per siswa, per tanggal, per kelas
2. Guru dapat melihat history absensi siswa (harian, per minggu, per bulan)
3. Sistem harus dapat export data dalam format CSV atau Excel
4. Data harus tersimpan aman dan tidak mudah dihapus

### FR5: Reporting & Analytics

1. Aplikasi harus support laporan dengan periode:
   - **Harian**: Data kehadiran siswa per hari per kelas
   - **Mingguan**: Total kehadiran/absen per minggu per siswa
   - **Bulanan**: Rekapitulasi kehadiran bulanan per siswa dengan presentase
   - **Semesteran**: Analisis kehadiran per semester dengan trend
   - **Custom Range**: User bisa pilih rentang tanggal custom (dari-sampai)
2. Laporan dapat difilter by:
   - **Kelas**: Single kelas atau multiple kelas
   - **Siswa**: Single siswa atau semua siswa di kelas
   - **Status**: Hadir, Izin, Sakit, Tanpa Keterangan, atau kombinasi
3. Laporan harus menampilkan:
   - Nama siswa, NISN, Kelas
   - Detail absensi per hari (tanggal + status)
   - Total: Jumlah Hadir, Izin, Sakit, Tanpa Keterangan
   - Presentase kehadiran
   - Trend/grafik (untuk laporan bulanan & semesteran)
4. Laporan dapat di-print dan di-export ke format CSV/Excel/PDF

### FR6: User Management

1. Aplikasi harus mendukung login untuk guru
2. Guru hanya bisa melihat data kelas yang mereka ajar
3. Admin sekolah dapat manage data guru dan kelas
4. Sistem harus manage user credentials dengan aman

### FR7: Student Database

1. Aplikasi harus memiliki database siswa dengan field minimum: NISN, Nama, Kelas, Foto
2. Admin dapat input/edit data siswa
3. Admin dapat generate QR code untuk siswa baru
4. Database dapat di-backup

---

## Non-Goals (Out of Scope)

1. **Integrasi dengan sistem pemerintah** - Aplikasi standalone, tidak perlu sync dengan database NISN pusat
2. **Notifikasi ke orang tua** - Fitur komunikasi ke orang tua tidak termasuk di fase ini
3. **Tracking lokasi siswa** - Aplikasi hanya mencatat kehadiran, bukan tracking movement
4. **Sistem punishment/reward** - Logika aturan absensi tidak termasuk (e.g., "absen 3x dikeluarkan")
5. **Multi-bahasa** - Aplikasi menggunakan bahasa Indonesia saja di fase awal
6. **Offline mode lengkap** - Fitur offline sync penuh tidak termasuk

---

## Design Considerations

### UI/UX Requirements

1. **Simple & Intuitive Interface** - Design harus mudah digunakan untuk semua umur (guru dan siswa)
2. **Large buttons & text** - Mempertimbangkan penggunaan oleh anak-anak SD
3. **Dark mode support** - Untuk usability di berbagai lighting condition
4. **Responsive Design** - Support smartphone (portrait/landscape) dan tablet

### User Flow - Scanning (Recommended: Automatic Approach)

```
Login → Buka Scanner → Scan QR → Sistem ekstrak NISN+Kelas+Tanggal dari QR →
Validasi Siswa → Cek duplikat → Terekam dengan timestamp → Tampilkan konfirmasi →
Daftar ter-scan (real-time) → Next scan atau View List atau Input Manual
```

- **Automatic (Recommended)**: Login → Scan → Auto-filled → Terekam
  - **Pros**: Lebih cepat, mengurangi human error, flow lebih streamlined
  - **Cons**: Requires QR code to contain kelas info

**Rekomendasi**: Gunakan **Automatic Approach** (QR code embed kelas+NISN) untuk efisiensi dan mengurangi error.

### User Flow - Manual Attendance Input

```
Login → Dashboard Guru → Pilih Kelas → Tab "Siswa Belum Absen" →
Pilih 1 atau Multiple Siswa → Set Status (Hadir/Izin/Sakit/Tanpa Keterangan) →
Tambah Catatan (optional) → Submit → Terekam & Audit Log
```

### User Flow - Reporting

```
Login → Dashboard → Menu "Laporan" → Pilih Tipe (Harian/Mingguan/Bulanan/Semesteran/Custom) →
Filter (Kelas/Siswa/Status) → View → Download/Export/Print
```

---

## Technical Considerations

### Technology Stack (Phase 1: Web Version)

**Frontend Web:**

- **Runtime/Package Manager**: Bun (JavaScript runtime & package manager, replaces Node.js)
- **Framework**: React 18+
- **Build Tool**: Vite
- **UI Library**: React Router v6 (routing)
- **State Management**: React Query
- **QR Code**:
  - Scan: `jsQR` atau `html5-qrcode`
  - Generate: `qrcode.react`
- **Styling**: Tailwind CSS dan shadcn/ui
- **HTTP Client**: Axios

**Backend API:**

- **Language**: Go (Golang)
- **Framework**: Gin Gonic (HTTP framework)
- **Dependency Injection**: Wire (code generation tool for dependency injection)
- **Database**: PostgreSQL + native SQL (database/sql package)
- **Authentication**: JWT (golang-jwt)
- **Middleware**: CORS, logging, error handling
- **File Upload**: untuk export CSV/Excel/PDF
- **Database Migration**: golang-migrate
- **Development**: Air (live reload)
- **Configuration**: Viper (config management)
- **Libraries**:
  - `uuid`: Generate unique ID
  - `bcrypt`: Password hashing
  - `pq`: PostgreSQL driver
  - `golang-jwt`: JWT authentication
  - `air`: Live reload untuk development
  - `golang-migrate`: Database migration tool
  - `viper`: Configuration & environment variables management

**Deployment & DevOps:**

- **Frontend Hosting**: Vercel, Netlify, atau self-hosted (Nginx)
- **Backend Hosting**: Docker + AWS/GCP/DigitalOcean atau self-hosted
- **Database**: PostgreSQL managed (AWS RDS) atau self-hosted
- **API Base URL**: REST API accessible via HTTP/HTTPS

**Additional:**

- **Authentication**: JWT token stored in localStorage (frontend) + HttpOnly cookies (optional)
- **CORS**: Configure backend untuk accept frontend origin
- **Rate Limiting**: Implement di Gin untuk prevent abuse

---

## Database Design

### ER Diagram (Text-based)

```
┌─────────────┐       ┌──────────────┐       ┌─────────────┐
│   users     │◄──────│   classes    │──────►│  students   │
│(guru/admin) │       │  (kelas)     │       │  (siswa)    │
└─────────────┘       └──────────────┘       └─────────────┘
       ▲                      ▲                      │
       │                      │                      │ 1:N
       │                      │                      ▼
       │                      │             ┌──────────────────┐
       └──────────────────────┴─────────────│   attendances    │
                     N:1                    │  (absensi)       │
                                            └──────────────────┘
                                                    │
                                            ┌──────────────────┐
                                            │attendance_audits │
                                            └──────────────────┘
```

### Table Specifications

#### 1. **users** (Guru & Admin)

```sql
CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  username VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  role ENUM('admin', 'guru') DEFAULT 'guru',
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  last_login TIMESTAMP
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
```

#### 2. **classes** (Kelas)

```sql
CREATE TABLE classes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  class_name VARCHAR(50) NOT NULL,  -- e.g., "1A", "1B", "2A"
  teacher_id UUID NOT NULL REFERENCES users(id),
  academic_year VARCHAR(9) NOT NULL,  -- e.g., "2024/2025"
  capacity INT DEFAULT 30,
  description TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(class_name, academic_year)
);

-- Indexes
CREATE INDEX idx_classes_teacher_id ON classes(teacher_id);
CREATE INDEX idx_classes_academic_year ON classes(academic_year);
```

#### 3. **students** (Siswa)

```sql
CREATE TABLE students (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  nisn VARCHAR(10) UNIQUE NOT NULL,
  full_name VARCHAR(150) NOT NULL,
  class_id UUID NOT NULL REFERENCES classes(id),
  date_of_birth DATE,
  gender ENUM('laki-laki', 'perempuan'),
  photo_url VARCHAR(500),
  qr_code_data TEXT NOT NULL,  -- JSON atau plain text: NISN|FullName|ClassName
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(nisn, class_id)  -- NISN per kelas (siswa bisa pindah kelas)
);

-- Indexes
CREATE INDEX idx_students_nisn ON students(nisn);
CREATE INDEX idx_students_class_id ON students(class_id);
CREATE INDEX idx_students_is_active ON students(is_active);
```

#### 4. **attendances** (Data Absensi)

```sql
CREATE TABLE attendances (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  student_id UUID NOT NULL REFERENCES students(id),
  class_id UUID NOT NULL REFERENCES classes(id),
  attendance_date DATE NOT NULL,
  status ENUM('hadir', 'izin', 'sakit', 'tanpa_keterangan') DEFAULT 'hadir',
  recorded_by UUID NOT NULL REFERENCES users(id),  -- guru yang input
  recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  scanned_at TIMESTAMP,  -- timestamp when QR was scanned (NULL if manual)
  notes TEXT,  -- keterangan tambahan
  is_manual BOOLEAN DEFAULT false,  -- true if input manual by teacher

  UNIQUE(student_id, class_id, attendance_date)  -- 1 absensi per siswa per hari per kelas
);

-- Indexes
CREATE INDEX idx_attendances_student_id ON attendances(student_id);
CREATE INDEX idx_attendances_class_id ON attendances(class_id);
CREATE INDEX idx_attendances_attendance_date ON attendances(attendance_date);
CREATE INDEX idx_attendances_status ON attendances(status);
CREATE INDEX idx_attendances_recorded_by ON attendances(recorded_by);
CREATE INDEX idx_attendances_composite ON attendances(class_id, attendance_date);
```

#### 5. **attendance_audits** (Audit Log - untuk track perubahan)

```sql
CREATE TABLE attendance_audits (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  attendance_id UUID NOT NULL REFERENCES attendances(id),
  old_status VARCHAR(50),  -- status sebelumnya
  new_status VARCHAR(50) NOT NULL,  -- status sesudahnya
  changed_by UUID NOT NULL REFERENCES users(id),  -- user yang ubah
  changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  reason TEXT  -- alasan perubahan
);

-- Indexes
CREATE INDEX idx_attendance_audits_attendance_id ON attendance_audits(attendance_id);
CREATE INDEX idx_attendance_audits_changed_by ON attendance_audits(changed_by);
CREATE INDEX idx_attendance_audits_changed_at ON attendance_audits(changed_at);
```

#### 6. **class_teachers** (Relasi Guru-Kelas, jika guru mengajar multiple kelas)

```sql
CREATE TABLE class_teachers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  teacher_id UUID NOT NULL REFERENCES users(id),
  class_id UUID NOT NULL REFERENCES classes(id),
  academic_year VARCHAR(9) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(teacher_id, class_id, academic_year)
);

-- Indexes
CREATE INDEX idx_class_teachers_teacher_id ON class_teachers(teacher_id);
CREATE INDEX idx_class_teachers_class_id ON class_teachers(class_id);
```

#### 7. **reports** (Cache untuk laporan agar lebih cepat)

```sql
CREATE TABLE reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  report_type ENUM('harian', 'mingguan', 'bulanan', 'semesteran') NOT NULL,
  class_id UUID NOT NULL REFERENCES classes(id),
  period_start DATE NOT NULL,
  period_end DATE NOT NULL,
  generated_by UUID NOT NULL REFERENCES users(id),
  generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  report_data JSONB,  -- cache report data

  UNIQUE(report_type, class_id, period_start, period_end)
);

-- Indexes
CREATE INDEX idx_reports_class_id ON reports(class_id);
CREATE INDEX idx_reports_report_type ON reports(report_type);
CREATE INDEX idx_reports_period ON reports(period_start, period_end);
```

### Key Design Decisions

1. **UUID Primary Keys** - Better for distributed systems & security (no sequential guessing)
2. **ENUM Types** - For status fields (hadir, izin, sakit, tanpa_keterangan)
3. **Composite Unique Constraint** - `attendances(student_id, class_id, attendance_date)` prevent duplikat absensi
4. **Audit Table** - Track semua perubahan status absensi dengan user trail
5. **JSONB for QR Data** - Flexible untuk future extensions (embed foto, barcode, dll)
6. **Indexes on Foreign Keys** - For fast JOIN queries on reporting
7. **Composite Index** - `(class_id, attendance_date)` untuk query laporan yang lebih cepat
8. **Soft Delete** - `is_active` flag instead of hard delete (audit trail)

### Performance Optimization

1. **Materialized Views** - Untuk reporting (pre-calculate stats)

   ```sql
   CREATE MATERIALIZED VIEW attendance_summary AS
   SELECT
     s.id,
     s.nisn,
     s.full_name,
     c.class_name,
     COUNT(CASE WHEN a.status = 'hadir' THEN 1 END) as hadir_count,
     COUNT(CASE WHEN a.status = 'izin' THEN 1 END) as izin_count,
     COUNT(CASE WHEN a.status = 'sakit' THEN 1 END) as sakit_count,
     COUNT(CASE WHEN a.status = 'tanpa_keterangan' THEN 1 END) as tanpa_count,
     COUNT(*) as total_records
   FROM students s
   JOIN classes c ON s.class_id = c.id
   LEFT JOIN attendances a ON s.id = a.student_id
   WHERE DATE_TRUNC('month', a.attendance_date) = DATE_TRUNC('month', CURRENT_DATE)
   GROUP BY s.id, c.id;
   ```

2. **Indexing Strategy**:
   - Index on frequently queried columns (class_id, student_id, attendance_date)
   - Composite index for common filter combinations
   - Index on status untuk filtering

3. **Query Caching** - Use reports table to cache generated reports

### Migration Strategy

```
001_init_schema.up.sql          -- Create all base tables
002_add_indexes.up.sql          -- Add all indexes
003_add_audit_table.up.sql      -- Add audit table
004_add_views.up.sql            -- Create materialized views
```

---

### Security Considerations

1. Validasi NISN harus strict (tidak bisa scan QR code orang lain)
2. Encrypt password dan sensitive data
3. Implement rate limiting untuk prevent brute force login
4. Audit log untuk setiap perubahan data absensi

### Performance Considerations

1. QR scan harus response time < 2 detik
2. Laporan harus load < 3 detik bahkan untuk data ribuan siswa
3. Cache data untuk reduce database load

### Data Privacy

1. Comply dengan regulasi perlindungan data anak
2. Hanya store data minimal yang diperlukan
3. Data harus dapat dihapus sesuai kebijakan sekolah

---

## API Documentation

### Base Configuration

- **Base URL**: `http://localhost:8080/api/v1` (development) atau `https://api.absensi.com/api/v1` (production)
- **Authentication**: JWT Bearer Token
- **Content-Type**: `application/json`
- **Response Format**: JSON

### Authentication

#### POST `/auth/login`

Login user (guru/admin)

**Request:**

```json
{
  "username": "guru_001",
  "password": "password123"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "guru_001",
      "full_name": "Ibu Nur Habibah",
      "email": "nur@school.edu",
      "role": "guru",
      "last_login": "2024-05-21T10:30:00Z"
    }
  },
  "message": "Login successful"
}
```

**Error (401 Unauthorized):**

```json
{
  "success": false,
  "error": "invalid_credentials",
  "message": "Username atau password salah"
}
```

#### POST `/auth/logout`

Logout user

**Headers:**

```
Authorization: Bearer {token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "message": "Logout successful"
}
```

#### POST `/auth/refresh-token`

Refresh JWT token

**Request:**

```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

---

### User Management

#### GET `/users/me`

Get current logged-in user info

**Headers:**

```
Authorization: Bearer {token}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "guru_001",
    "full_name": "Ibu Nur Habibah",
    "email": "nur@school.edu",
    "role": "guru",
    "is_active": true,
    "created_at": "2024-01-15T08:00:00Z"
  }
}
```

#### GET `/users`

List all users (admin only)

**Query Parameters:**

- `role`: "guru" atau "admin" (optional)
- `is_active`: true/false (optional)
- `page`: 1 (default)
- `limit`: 10 (default)

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "guru_001",
      "full_name": "Ibu Nur Habibah",
      "email": "nur@school.edu",
      "role": "guru",
      "is_active": true,
      "created_at": "2024-01-15T08:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

#### POST `/users`

Create new user (admin only)

**Request:**

```json
{
  "username": "guru_002",
  "email": "guru2@school.edu",
  "password": "securepass123",
  "full_name": "Pak Bambang Irawan",
  "role": "guru"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": "650e8400-e29b-41d4-a716-446655440001",
    "username": "guru_002",
    "full_name": "Pak Bambang Irawan",
    "email": "guru2@school.edu",
    "role": "guru",
    "is_active": true,
    "created_at": "2024-05-21T10:00:00Z"
  },
  "message": "User created successfully"
}
```

#### PUT `/users/{user_id}`

Update user (admin or self)

**Request:**

```json
{
  "full_name": "Ibu Nur Habibah S.Pd",
  "email": "nur.habibah@school.edu"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "guru_001",
    "full_name": "Ibu Nur Habibah S.Pd",
    "email": "nur.habibah@school.edu",
    "role": "guru"
  },
  "message": "User updated successfully"
}
```

---

### Class Management

#### GET `/classes`

Get list of classes (teacher only sees their own classes)

**Query Parameters:**

- `academic_year`: "2024/2025" (optional)
- `page`: 1 (default)
- `limit`: 20 (default)

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "750e8400-e29b-41d4-a716-446655440002",
      "class_name": "1A",
      "teacher_id": "550e8400-e29b-41d4-a716-446655440000",
      "teacher_name": "Ibu Nur Habibah",
      "academic_year": "2024/2025",
      "capacity": 30,
      "student_count": 28,
      "description": "Kelas 1 Angkatan A",
      "created_at": "2024-01-10T09:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 12,
    "total_pages": 1
  }
}
```

#### GET `/classes/{class_id}`

Get class detail with list of students

**Query Parameters:**

- `include_attendance`: true/false (optional, for today's attendance)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440000",
    "academic_year": "2024/2025",
    "capacity": 30,
    "students": [
      {
        "id": "850e8400-e29b-41d4-a716-446655440003",
        "nisn": "1234567890",
        "full_name": "Ahmad Rizki Pratama",
        "qr_code_data": "1234567890|Ahmad Rizki Pratama|1A",
        "attendance_today": "hadir",
        "scanned_at": "2024-05-21T07:05:00Z"
      }
    ]
  }
}
```

#### POST `/classes`

Create new class (admin only)

**Request:**

```json
{
  "class_name": "2B",
  "teacher_id": "550e8400-e29b-41d4-a716-446655440000",
  "academic_year": "2024/2025",
  "capacity": 32,
  "description": "Kelas 2 Angkatan B"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440004",
    "class_name": "2B",
    "teacher_id": "550e8400-e29b-41d4-a716-446655440000",
    "academic_year": "2024/2025",
    "capacity": 32,
    "description": "Kelas 2 Angkatan B",
    "created_at": "2024-05-21T10:30:00Z"
  }
}
```

---

### Student Management

#### GET `/students`

Get list of students (with optional filters)

**Query Parameters:**

- `class_id`: UUID (optional)
- `is_active`: true/false (optional)
- `search`: "nama atau NISN" (optional)
- `page`: 1 (default)
- `limit`: 20 (default)

**Response (200 OK):**

```json
{
  "success": true,
  "data": [
    {
      "id": "850e8400-e29b-41d4-a716-446655440003",
      "nisn": "1234567890",
      "full_name": "Ahmad Rizki Pratama",
      "class_id": "750e8400-e29b-41d4-a716-446655440002",
      "class_name": "1A",
      "date_of_birth": "2018-03-15",
      "gender": "laki-laki",
      "photo_url": "https://api.absensi.com/photos/850e8400.jpg",
      "is_active": true,
      "created_at": "2024-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 450,
    "total_pages": 23
  }
}
```

#### GET `/students/{student_id}`

Get single student detail

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "850e8400-e29b-41d4-a716-446655440003",
    "nisn": "1234567890",
    "full_name": "Ahmad Rizki Pratama",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "date_of_birth": "2018-03-15",
    "gender": "laki-laki",
    "photo_url": "https://api.absensi.com/photos/850e8400.jpg",
    "qr_code_data": "1234567890|Ahmad Rizki Pratama|1A",
    "is_active": true,
    "created_at": "2024-01-20T10:00:00Z"
  }
}
```

#### POST `/students`

Create new student (admin only)

**Request:**

```json
{
  "nisn": "1234567891",
  "full_name": "Siti Nur Azizah",
  "class_id": "750e8400-e29b-41d4-a716-446655440002",
  "date_of_birth": "2018-05-22",
  "gender": "perempuan",
  "photo_url": "https://example.com/photo.jpg"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": "950e8400-e29b-41d4-a716-446655440005",
    "nisn": "1234567891",
    "full_name": "Siti Nur Azizah",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "qr_code_data": "1234567891|Siti Nur Azizah|1A",
    "is_active": true,
    "created_at": "2024-05-21T11:00:00Z"
  },
  "message": "Student created successfully"
}
```

#### GET `/students/{student_id}/qrcode`

Generate/download QR code for a student

**Response:**

```
[Binary image file - PNG format]
Content-Type: image/png
```

#### PUT `/students/{student_id}`

Update student info (admin only)

**Request:**

```json
{
  "full_name": "Siti Nur Azizah Putri",
  "photo_url": "https://example.com/photo-updated.jpg"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "950e8400-e29b-41d4-a716-446655440005",
    "nisn": "1234567891",
    "full_name": "Siti Nur Azizah Putri",
    "class_id": "750e8400-e29b-41d4-a716-446655440002"
  },
  "message": "Student updated successfully"
}
```

---

### Attendance - QR Code Scanning

#### POST `/attendances/qr-scan`

Record attendance via QR code scan

**Request:**

```json
{
  "qr_code_data": "1234567890|Ahmad Rizki Pratama|1A",
  "class_id": "750e8400-e29b-41d4-a716-446655440002",
  "scanned_at": "2024-05-21T07:05:00Z"
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "id": "a50e8400-e29b-41d4-a716-446655440006",
    "student_id": "850e8400-e29b-41d4-a716-446655440003",
    "student_name": "Ahmad Rizki Pratama",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "attendance_date": "2024-05-21",
    "status": "hadir",
    "scanned_at": "2024-05-21T07:05:00Z",
    "is_manual": false,
    "recorded_by": "550e8400-e29b-41d4-a716-446655440000"
  },
  "message": "Attendance recorded successfully"
}
```

**Error (409 Conflict - Duplicate scan):**

```json
{
  "success": false,
  "error": "duplicate_scan",
  "message": "Siswa Ahmad Rizki Pratama sudah ter-scan hari ini"
}
```

**Error (404 Not Found):**

```json
{
  "success": false,
  "error": "student_not_found",
  "message": "QR code tidak valid atau siswa tidak ditemukan"
}
```

---

### Attendance - Manual Input

#### POST `/attendances/manual`

Record attendance manually by teacher

**Request:**

```json
{
  "class_id": "750e8400-e29b-41d4-a716-446655440002",
  "attendance_date": "2024-05-21",
  "students": [
    {
      "student_id": "850e8400-e29b-41d4-a716-446655440003",
      "status": "hadir"
    },
    {
      "student_id": "950e8400-e29b-41d4-a716-446655440005",
      "status": "izin",
      "notes": "Sakit gigi"
    },
    {
      "student_id": "b50e8400-e29b-41d4-a716-446655440007",
      "status": "sakit"
    }
  ]
}
```

**Response (201 Created):**

```json
{
  "success": true,
  "data": {
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "attendance_date": "2024-05-21",
    "records_created": 3,
    "records": [
      {
        "id": "a50e8400-e29b-41d4-a716-446655440006",
        "student_id": "850e8400-e29b-41d4-a716-446655440003",
        "status": "hadir",
        "is_manual": true,
        "recorded_at": "2024-05-21T08:30:00Z"
      }
    ]
  },
  "message": "Attendance recorded successfully"
}
```

#### GET `/attendances/{class_id}/{date}`

Get attendance for a class on specific date

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "attendance_date": "2024-05-21",
    "total_students": 28,
    "summary": {
      "hadir": 26,
      "izin": 1,
      "sakit": 1,
      "tanpa_keterangan": 0
    },
    "records": [
      {
        "id": "a50e8400-e29b-41d4-a716-446655440006",
        "student_id": "850e8400-e29b-41d4-a716-446655440003",
        "student_name": "Ahmad Rizki Pratama",
        "nisn": "1234567890",
        "status": "hadir",
        "scanned_at": "2024-05-21T07:05:00Z",
        "is_manual": false,
        "recorded_by": "550e8400-e29b-41d4-a716-446655440000"
      }
    ]
  }
}
```

#### PUT `/attendances/{attendance_id}`

Update attendance record (admin/teacher)

**Request:**

```json
{
  "status": "izin",
  "notes": "Izin untuk ke dokter",
  "reason": "Update karena siswa membawa surat izin"
}
```

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "id": "a50e8400-e29b-41d4-a716-446655440006",
    "student_id": "850e8400-e29b-41d4-a716-446655440003",
    "status": "izin",
    "notes": "Izin untuk ke dokter",
    "updated_at": "2024-05-21T10:00:00Z"
  },
  "message": "Attendance updated successfully",
  "audit_log": {
    "old_status": "hadir",
    "new_status": "izin",
    "changed_by": "550e8400-e29b-41d4-a716-446655440000",
    "changed_at": "2024-05-21T10:00:00Z",
    "reason": "Update karena siswa membawa surat izin"
  }
}
```

---

### Reporting

#### GET `/reports/daily`

Get daily attendance report

**Query Parameters:**

- `class_id`: UUID (required)
- `date`: "2024-05-21" (required)
- `format`: "json" atau "csv" atau "pdf" (optional, default: json)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "report_type": "harian",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "date": "2024-05-21",
    "total_students": 28,
    "summary": {
      "hadir": 26,
      "izin": 1,
      "sakit": 1,
      "tanpa_keterangan": 0,
      "hadir_percentage": 92.86
    },
    "records": [
      {
        "nisn": "1234567890",
        "student_name": "Ahmad Rizki Pratama",
        "status": "hadir",
        "scanned_at": "2024-05-21T07:05:00Z"
      }
    ]
  },
  "generated_at": "2024-05-21T10:30:00Z"
}
```

#### GET `/reports/monthly`

Get monthly attendance report

**Query Parameters:**

- `class_id`: UUID (required)
- `month`: "2024-05" (required, format: YYYY-MM)
- `format`: "json" atau "csv" atau "pdf" (optional)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "report_type": "bulanan",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "period": "May 2024",
    "total_days": 22,
    "summary": {
      "total_students": 28,
      "avg_hadir_percentage": 95.2,
      "total_izin": 8,
      "total_sakit": 12,
      "total_tanpa_keterangan": 5
    },
    "student_stats": [
      {
        "nisn": "1234567890",
        "student_name": "Ahmad Rizki Pratama",
        "hadir": 21,
        "izin": 0,
        "sakit": 1,
        "tanpa_keterangan": 0,
        "attendance_percentage": 95.45
      }
    ]
  },
  "generated_at": "2024-05-21T10:30:00Z"
}
```

#### GET `/reports/semester`

Get semester attendance report

**Query Parameters:**

- `class_id`: UUID (required)
- `semester`: "1" atau "2" (required)
- `academic_year`: "2024/2025" (required)
- `format`: "json" atau "csv" atau "pdf" (optional)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "report_type": "semesteran",
    "class_id": "750e8400-e29b-41d4-a716-446655440002",
    "class_name": "1A",
    "period": "Semester 1 - 2024/2025",
    "duration_days": 87,
    "summary": {
      "avg_attendance": 94.1,
      "total_izin": 45,
      "total_sakit": 52,
      "total_tanpa_keterangan": 18
    },
    "trend": [
      {
        "month": "July 2024",
        "attendance_percentage": 96.2
      },
      {
        "month": "August 2024",
        "attendance_percentage": 94.8
      }
    ]
  },
  "generated_at": "2024-05-21T10:30:00Z"
}
```

#### GET `/reports/student/{student_id}`

Get attendance history for specific student

**Query Parameters:**

- `from_date`: "2024-05-01" (optional)
- `to_date`: "2024-05-31" (optional)

**Response (200 OK):**

```json
{
  "success": true,
  "data": {
    "student_id": "850e8400-e29b-41d4-a716-446655440003",
    "nisn": "1234567890",
    "student_name": "Ahmad Rizki Pratama",
    "class_name": "1A",
    "period": "May 1 - May 31, 2024",
    "summary": {
      "total_days": 22,
      "hadir": 21,
      "izin": 0,
      "sakit": 1,
      "tanpa_keterangan": 0,
      "attendance_percentage": 95.45
    },
    "daily_records": [
      {
        "date": "2024-05-21",
        "status": "hadir",
        "scanned_at": "2024-05-21T07:05:00Z"
      }
    ]
  }
}
```

#### POST `/reports/export`

Export report to CSV/Excel/PDF

**Request:**

```json
{
  "report_type": "monthly",
  "class_id": "750e8400-e29b-41d4-a716-446655440002",
  "month": "2024-05",
  "format": "excel"
}
```

**Response:**

```
[Binary file - .xlsx or .csv or .pdf format]
Content-Type: application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
Content-Disposition: attachment; filename="Laporan_Absensi_1A_May2024.xlsx"
```

---

### Error Handling

All errors follow this standard format:

```json
{
  "success": false,
  "error": "error_code",
  "message": "Human readable error message",
  "details": {} // Optional: additional error details
}
```

**Common Error Codes:**

- `400_bad_request` - Invalid request parameters
- `401_unauthorized` - Authentication required or invalid token
- `403_forbidden` - User doesn't have permission
- `404_not_found` - Resource not found
- `409_conflict` - Duplicate entry or conflict
- `422_unprocessable_entity` - Validation error
- `500_internal_error` - Server error
- `503_service_unavailable` - Service maintenance

**Example Error Response (422 Unprocessable Entity):**

```json
{
  "success": false,
  "error": "validation_error",
  "message": "Validation failed",
  "details": {
    "email": "Email format is invalid",
    "full_name": "Full name is required"
  }
}
```

---

### Rate Limiting

- **Limit**: 100 requests per minute per user
- **Headers**:
  ```
  X-RateLimit-Limit: 100
  X-RateLimit-Remaining: 95
  X-RateLimit-Reset: 1716345600
  ```
- **Error (429 Too Many Requests)**:
  ```json
  {
    "success": false,
    "error": "rate_limit_exceeded",
    "message": "Terlalu banyak request. Coba lagi dalam 60 detik"
  }
  ```

---

### Pagination

All list endpoints support pagination with these query parameters:

- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)
- `sort_by`: Column name (default: created_at)
- `sort_order`: "asc" atau "desc" (default: desc)

**Response format includes:**

```json
{
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150,
    "total_pages": 8
  }
}
```

1. **Kecepatan Absensi** - Rata-rata waktu absensi per siswa < 3 detik per scan
2. **Akurasi Data** - 99% data absensi terekam dengan benar tanpa duplikat
3. **User Adoption** - 80%+ guru aktif menggunakan aplikasi dalam 1 bulan pertama
4. **Error Rate** - Fail rate scan < 1%
5. **Data Completeness** - 95%+ data absensi siswa terkumpul setiap hari
6. **Teacher Satisfaction** - NPS score > 7 dari feedback guru

---

## Open Questions untuk Klarifikasi Lanjutan

1. **Offline Functionality** - Phase 2 feature: Bisakah scan QR code offline + auto-sync ketika online?
2. **QR Code Distribution** - QR code dicetak di kartu siswa atau digital di phone?
3. **Late Check-in** - Adakah tolerance time untuk status "terlambat"?
4. **Admin Dashboard** - Perlu dashboard untuk admin sekolah (oversight semua kelas)?
5. **Mobile App** - Kapan direncanakan versi mobile app (React Native/Flutter)? Phase 2?
6. **QR Code Data** - QR code include: NISN (10 digit) + Nama + Kelas, atau tambah Foto/Barcode?
7. **Encryption** - Apakah QR code harus encrypted atau plain text NISN?
8. **Performance Target** - Berapa jumlah siswa di satu kelas? (affect database design)

---

---

## Implementation Roadmap

### Phase 1: MVP (4-6 weeks)

- ✅ QR Code generation & scanning (web)
- ✅ Basic attendance recording (scan + manual input)
- ✅ Laporan harian per kelas
- ✅ User authentication (guru login)
- ✅ Basic admin (manage siswa & generate QR)

### Phase 2: Enhancement (4 weeks)

- 📋 Laporan lanjutan (mingguan/bulanan/semesteran)
- 📋 Offline scanning + auto-sync
- 📋 Mobile app (React Native)
- 📋 Admin dashboard (oversight semua kelas)
- 📋 Notification system (email to teachers)

### Phase 3: Advanced (Future)

- 📋 Integrasi SISFO sekolah
- 📋 Parent portal (notifikasi ke orang tua)
- 📋 Analytics & AI insights
- 📋 Multi-sekolah support

---
