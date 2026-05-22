# Task Implementasi: Penyesuaian API Report dengan Frontend

Dokumen ini adalah tugas (*task*) untuk memperbaiki/menyesuaikan endpoint API *Reports* agar selaras dengan kode *React Query* di *Frontend*.

## Konteks Kode Frontend yang Berjalan

Frontend menggunakan Axios (diasumsikan `api.get` secara otomatis menambahkan *base URL* seperti `/api/v1`) dan mendefinisikan *hooks* berikut:

```typescript
import { useQuery } from "@tanstack/react-query";
import api from "@/utils/api";
import type { 
  DailyReport, 
  MonthlyReport, 
  SemesterReport,
  ReportQueryParams 
} from "@/types";

/**
 * Fetch daily report
 */
export const useDailyReport = (params: ReportQueryParams) => {
  return useQuery({
    queryKey: ["report", "daily", params.class_id, params.date],
    queryFn: async () => {
      if (!params.class_id || !params.date) return null;
      const res = await api.get<{ success: boolean; data: DailyReport }>(
        `/reports/daily?class_id=${params.class_id}&date=${params.date}`
      );
      return res.data.data;
    },
    enabled: !!params.class_id && !!params.date,
  });
};

/**
 * Fetch monthly report
 */
export const useMonthlyReport = (params: ReportQueryParams) => {
  return useQuery({
    queryKey: ["report", "monthly", params.class_id, params.month],
    queryFn: async () => {
      if (!params.class_id || !params.month) return null;
      const res = await api.get<{ success: boolean; data: MonthlyReport }>(
        `/reports/monthly?class_id=${params.class_id}&month=${params.month}`
      );
      return res.data.data;
    },
    enabled: !!params.class_id && !!params.month,
  });
};

/**
 * Fetch semester report
 */
export const useSemesterReport = (params: ReportQueryParams) => {
  return useQuery({
    queryKey: ["report", "semester", params.class_id, params.semester, params.academic_year],
    queryFn: async () => {
      if (!params.class_id || !params.semester || !params.academic_year) return null;
      const res = await api.get<{ success: boolean; data: SemesterReport }>(
        `/reports/semester?class_id=${params.class_id}&semester=${params.semester}&academic_year=${params.academic_year}`
      );
      return res.data.data;
    },
    enabled: !!params.class_id && !!params.semester && !!params.academic_year,
  });
};
```

## Tugas yang Harus Dijalankan

1. **Analisis Struktur JSON Response**:
   Saat ini, API di sisi *Backend* telah dibuat dan merespons dengan JSON berstruktur:
   `{ "success": true, "data": { ... } }`.
   *   **Tugas**: Pastikan bahwa properti (kunci JSON) yang dikembalikan di dalam `data` dari backend benar-benar memiliki field yang diharapkan oleh tipe TypeScript `DailyReport`, `MonthlyReport`, dan `SemesterReport` di Frontend.
   *   *(Catatan untuk AI: Karena Anda mungkin tidak memiliki file `types.ts` dari frontend, Anda perlu memberikan struktur respons API Backend saat ini kepada *User* dan bertanya "Apakah struktur JSON ini sudah sesuai dengan `types.ts` di Frontend Anda?", lalu perbarui `domain/report.go` jika ada ketidaksesuaian penamaan field (misal: camelCase vs snake_case).)*

2. **Periksa Parameter URL**:
   *   Frontend akan mengirimkan URL dengan query seperti: `academic_year=2024/2025`
   *   Karakter `/` dalam query parameter bisa jadi akan di-_encode_ secara otomatis oleh browser menjadi `%2F` (misal: `academic_year=2024%2F2025`).
   *   **Tugas**: Periksa logika di `internal/handler/report_handler.go` atau `internal/service/report_service.go` untuk memastikan backend tidak menolak nilai `academic_year` yang di-_encode_ dengan URL (bila framework Gin tidak men-decode-nya dengan sempurna).

3. **CORS Preflight Check (Opsional/Pastikan Saja)**:
   *   Jika *Frontend* berjalan di port yang berbeda dan mengalami error jaringan saat mengeksekusi `api.get`, pastikan konfigurasi `FRONTEND_URL` di dalam `docker-compose.yml` backend telah merujuk ke URL Frontend yang tepat (contoh: `http://localhost:3000`).

## Langkah Eksekusi (Action Items untuk Developer/AI Baru):
- Baca `internal/domain/report.go` untuk mengetahui struktur respons JSON saat ini.
- Diskusikan dengan *User* untuk mendapatkan definisi interface `DailyReport`, `MonthlyReport`, `SemesterReport` dari TypeScript Frontend.
- Lakukan pembaruan (refactoring) nama-nama *json tags* di `domain/report.go` jika terdapat perbedaan kapitalisasi atau susunan field.
- Jalankan ulang (rebuild) backend jika terdapat perubahan kode.
