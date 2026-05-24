package handler

import (
	"errors"
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
			c.JSON(http.StatusConflict, response.Error("Data yang dimasukkan sudah terdaftar di sistem. Silakan periksa kembali."))
			return
		case "23503": // foreign_key_violation
			c.JSON(http.StatusBadRequest, response.Error("Data referensi yang dipilih tidak valid atau tidak ditemukan."))
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
