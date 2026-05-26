package handler

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"alhikmah-attendance-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func (h *ReportHandler) exportDaily(c *gin.Context, report *domain.DailyReport, format string) {
	filename := fmt.Sprintf("Laporan Kehadiran Harian Kelas %s", report.ClassName)

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		writer.Write([]string{"NISN", "Nama Siswa", "Status", "Waktu Absen", "Metode"})
		for _, rec := range report.Records {
			scanned := "-"
			if rec.ScannedAt != nil {
				scanned = *rec.ScannedAt
			}
			method := "QR Scan"
			if rec.IsManual {
				method = "Manual"
			}
			writer.Write([]string{rec.NISN, rec.StudentName, rec.Status, scanned, method})
		}
		return
	}

	if format == "excel" {
		f := excelize.NewFile()
		defer f.Close()
		sheet := report.Date
		if sheet == "" {
			sheet = "Harian"
		}
		f.SetSheetName("Sheet1", sheet)

		headers := []string{"NISN", "Nama Siswa", "Status", "Waktu Absen", "Metode"}
		for i, hd := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, hd)
		}

		for i, rec := range report.Records {
			row := i + 2
			scanned := "-"
			if rec.ScannedAt != nil {
				scanned = *rec.ScannedAt
			}
			method := "QR Scan"
			if rec.IsManual {
				method = "Manual"
			}
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec.NISN)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), rec.StudentName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), rec.Status)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), scanned)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), method)
		}

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
		f.Write(c.Writer)
		return
	}
}

func (h *ReportHandler) exportMonthly(c *gin.Context, report *domain.MonthlyReport, format string) {
	filename := fmt.Sprintf("Laporan Kehadiran Bulanan Kelas %s", report.ClassName)

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		writer.Write([]string{"NISN", "Nama Siswa", "Hadir", "Izin", "Sakit", "Alpa", "Persentase Kehadiran (%)"})
		for _, rec := range report.StudentStats {
			writer.Write([]string{
				rec.NISN,
				rec.StudentName,
				strconv.Itoa(rec.Hadir),
				strconv.Itoa(rec.Izin),
				strconv.Itoa(rec.Sakit),
				strconv.Itoa(rec.TanpaKeterangan),
				fmt.Sprintf("%.2f", rec.AttendancePercentage),
			})
		}
		return
	}

	if format == "excel" {
		f := excelize.NewFile()
		defer f.Close()
		sheet := report.Period
		if sheet == "" {
			sheet = "Bulanan"
		}
		f.SetSheetName("Sheet1", sheet)

		headers := []string{"NISN", "Nama Siswa", "Hadir", "Izin", "Sakit", "Alpa", "Persentase Kehadiran (%)"}
		for i, hd := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, hd)
		}

		for i, rec := range report.StudentStats {
			row := i + 2
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec.NISN)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), rec.StudentName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), rec.Hadir)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), rec.Izin)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rec.Sakit)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rec.TanpaKeterangan)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%.2f", rec.AttendancePercentage))
		}

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
		f.Write(c.Writer)
		return
	}
}

func (h *ReportHandler) exportSemester(c *gin.Context, report *domain.SemesterReport, format string) {
	filename := fmt.Sprintf("Laporan Kehadiran Semester Kelas %s", report.ClassName)

	if format == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", filename))

		writer := csv.NewWriter(c.Writer)
		defer writer.Flush()

		writer.Write([]string{"NISN", "Nama Siswa", "Hadir", "Izin", "Sakit", "Alpa", "Persentase Kehadiran (%)"})
		for _, rec := range report.StudentStats {
			writer.Write([]string{
				rec.NISN,
				rec.StudentName,
				strconv.Itoa(rec.Hadir),
				strconv.Itoa(rec.Izin),
				strconv.Itoa(rec.Sakit),
				strconv.Itoa(rec.TanpaKeterangan),
				fmt.Sprintf("%.2f", rec.AttendancePercentage),
			})
		}
		return
	}

	if format == "excel" {
		f := excelize.NewFile()
		defer f.Close()
		
		// Excel limits sheet name to 31 chars and disallows some characters like /
		sheet := report.Period
		sheet = strings.ReplaceAll(sheet, "/", "-")
		if len(sheet) > 31 {
			sheet = sheet[:31] 
		}
		if sheet == "" {
			sheet = "Semesteran"
		}
		f.SetSheetName("Sheet1", sheet)

		headers := []string{"NISN", "Nama Siswa", "Hadir", "Izin", "Sakit", "Alpa", "Persentase Kehadiran (%)"}
		for i, hd := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, hd)
		}

		for i, rec := range report.StudentStats {
			row := i + 2
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), rec.NISN)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), rec.StudentName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), rec.Hadir)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), rec.Izin)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), rec.Sakit)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), rec.TanpaKeterangan)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), fmt.Sprintf("%.2f", rec.AttendancePercentage))
		}

		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.xlsx", filename))
		f.Write(c.Writer)
		return
	}
}
