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

func HandleAppError(c *gin.Context, err error, defaultMsg string) {
	fmt.Printf("[ERROR] %s: %v\n", defaultMsg, err)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505": // unique_violation
			msg := "The data already exists in the system."
			if pqErr.Detail != "" {
				if strings.HasPrefix(pqErr.Detail, "Key (") {
					parts := strings.Split(pqErr.Detail, ")=")
					if len(parts) >= 2 {
						fieldPart := strings.TrimPrefix(parts[0], "Key (")
						valPart := strings.Split(parts[1], ")")[0]
						valPart = strings.TrimPrefix(valPart, "(")

						// Field name mapping
						fieldMap := map[string]string{
							"nisn":                      "NISN",
							"username":                  "Username",
							"email":                     "Email",
							"class_name, academic_year": "Class name for this academic year",
							"teacher_id, class_id":      "Teacher assignment in this class",
						}

						fieldName := fieldPart
						if mapped, exists := fieldMap[fieldPart]; exists {
							fieldName = mapped
						}

						msg = fmt.Sprintf("%s '%s' already exists in the system. Please use a different value.", fieldName, valPart)
					}
				}
			}
			c.JSON(http.StatusConflict, response.Error(msg))
			return
		case "23503": // foreign_key_violation
			msg := "The referenced data is invalid or does not exist."
			if pqErr.Detail != "" {
				msg = msg + " Detail: " + pqErr.Detail
			}
			c.JSON(http.StatusBadRequest, response.Error(msg))
			return
		case "22P02": // invalid_text_representation (e.g. enum invalid)
			c.JSON(http.StatusBadRequest, response.Error("Invalid data format for the specified field type. Detail: "+pqErr.Message))
			return
		}
	}

	if strings.Contains(err.Error(), "not found") {
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "missing required") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "is required") ||
		strings.Contains(errMsg, "must be") ||
		strings.Contains(errMsg, "cannot be empty") ||
		strings.Contains(errMsg, "cannot delete") ||
		strings.Contains(errMsg, "already") ||
		strings.Contains(errMsg, "does not have access") {
		c.JSON(http.StatusBadRequest, response.Error(errMsg))
		return
	}

	if defaultMsg != "" {
		c.JSON(http.StatusInternalServerError, response.Error(fmt.Sprintf("%s: %s", defaultMsg, errMsg)))
	} else {
		c.JSON(http.StatusInternalServerError, response.Error(fmt.Sprintf("Internal Server Error: %s", errMsg)))
	}
}
