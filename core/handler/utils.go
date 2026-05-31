package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"alhikmah-attendance-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func handleDBError(c *gin.Context, err error) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			msg := "Data yang dimasukkan sudah terdaftar di sistem."
			if pqErr.Detail != "" {
				if strings.HasPrefix(pqErr.Detail, "Key (") {
					parts := strings.Split(pqErr.Detail, ")=")
					if len(parts) >= 2 {
						fieldPart := strings.TrimPrefix(parts[0], "Key (")
						valPart := strings.Split(parts[1], ")")[0]
						valPart = strings.TrimPrefix(valPart, "(")

						// Mapping nama field ke bahasa Indonesia
						fieldMap := map[string]string{
							"nisn":                        "NISN",
							"username":                    "Username",
							"email":                       "Email",
							"class_name, academic_year":   "Nama Kelas pada Tahun Ajaran tersebut",
							"teacher_id, class_id":        "Guru pada Kelas tersebut",
						}

						fieldName := fieldPart
						if mapped, exists := fieldMap[fieldPart]; exists {
							fieldName = mapped
						}

						msg = fmt.Sprintf("%s '%s' sudah terdaftar di sistem. Silakan gunakan data lain.", fieldName, valPart)
					}
				}
			}
			c.JSON(http.StatusConflict, response.Error(msg))
			return
		case "23503": // foreign_key_violation
			msg := "Data referensi yang dipilih tidak valid atau tidak ditemukan."
			if pqErr.Detail != "" {
				msg = msg + " Detail: " + pqErr.Detail
			}
			c.JSON(http.StatusBadRequest, response.Error(msg))
			return
		case "22P02": // invalid_text_representation (e.g. enum invalid)
			c.JSON(http.StatusBadRequest, response.Error("Format data tidak valid sesuai tipe yang diizinkan sistem. Detail: "+pqErr.Message))
			return
		}
	}

	if strings.Contains(err.Error(), "not found") {
		// Use the service's error message but standard format
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusInternalServerError, response.Error("Terjadi kesalahan pada server. Silakan coba beberapa saat lagi."))
}
