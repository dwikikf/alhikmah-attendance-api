package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"alhikmah-attendance-api/core/handler"
	"alhikmah-attendance-api/core/mocks"

	"github.com/stretchr/testify/assert"
)

func TestScanQR_MissingNISN(t *testing.T) {
	mockAttService := new(mocks.AttendanceService)
	h := handler.NewAttendanceHandler(mockAttService)
	router := newTestRouter()
	router.POST("/attendance/scan", h.ScanQR)

	w := httptest.NewRecorder()
	// Missing NISN in body
	req, _ := http.NewRequest("POST", "/attendance/scan", bytes.NewBuffer([]byte(`{}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestManualAttendance_InvalidStatus(t *testing.T) {
	mockAttService := new(mocks.AttendanceService)
	h := handler.NewAttendanceHandler(mockAttService)
	router := newTestRouter()
	router.POST("/attendance/manual", h.ManualInput)

	w := httptest.NewRecorder()
	// Invalid status
	req, _ := http.NewRequest("POST", "/attendance/manual", bytes.NewBuffer([]byte(`{"student_id":"123", "status":"invalid"}`)))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
