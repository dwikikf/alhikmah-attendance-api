package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"alhikmah-attendance-api/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestHandleAppError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		err          error
		defaultMsg   string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Generic Error",
			err:          errors.New("some random error"),
			defaultMsg:   "Failed to process request",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"success":false,"message":"Failed to process request: some random error"}`,
		},
		{
			name:         "Generic Error without Default Msg",
			err:          errors.New("some random error"),
			defaultMsg:   "",
			expectedCode: http.StatusInternalServerError,
			expectedBody: `{"success":false,"message":"Internal Server Error: some random error"}`,
		},
		{
			name:         "Not Found Error",
			err:          errors.New("student not found"),
			defaultMsg:   "",
			expectedCode: http.StatusNotFound,
			expectedBody: `{"success":false,"message":"student not found"}`,
		},
		{
			name: "Unique Violation Error (NISN)",
			err: &pq.Error{
				Code:   "23505",
				Detail: "Key (nisn)=(1234567890) already exists.",
			},
			defaultMsg:   "Failed to create",
			expectedCode: http.StatusConflict,
			expectedBody: `{"success":false,"message":"NISN '1234567890' already exists in the system. Please use a different value."}`,
		},
		{
			name: "Foreign Key Violation Error",
			err: &pq.Error{
				Code:   "23503",
				Detail: "Key (class_id)=(xyz) is not present in table \"classes\".",
			},
			defaultMsg:   "Failed to create",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"success":false,"message":"The referenced data is invalid or does not exist. Detail: Key (class_id)=(xyz) is not present in table \"classes\"."}`,
		},
		{
			name: "Invalid Text Representation (Enum)",
			err: &pq.Error{
				Code:    "22P02",
				Message: "invalid input value for enum",
			},
			defaultMsg:   "",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"success":false,"message":"Invalid data format for the specified field type. Detail: invalid input value for enum"}`,
		},
		{
			name:         "Validation Error (missing fields)",
			err:          errors.New("missing required fields"),
			defaultMsg:   "Failed to create",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"success":false,"message":"missing required fields"}`,
		},
		{
			name:         "Validation Error (invalid role)",
			err:          errors.New("invalid role"),
			defaultMsg:   "Failed to update",
			expectedCode: http.StatusBadRequest,
			expectedBody: `{"success":false,"message":"invalid role"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler.HandleAppError(c, tt.err, tt.defaultMsg)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
