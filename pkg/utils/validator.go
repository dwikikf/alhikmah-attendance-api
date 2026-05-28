package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError converts validator.ValidationErrors into a map of readable string messages.
func FormatValidationError(err error) map[string]string {
	errs := make(map[string]string)
	
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			// Convert Field name to snake_case or lowercase
			field := strings.ToLower(e.Field())
			
			switch e.Tag() {
			case "required":
				errs[field] = fmt.Sprintf("Kolom %s wajib diisi", field)
			case "min":
				errs[field] = fmt.Sprintf("Kolom %s minimal %s karakter", field, e.Param())
			case "max":
				errs[field] = fmt.Sprintf("Kolom %s maksimal %s karakter", field, e.Param())
			case "email":
				errs[field] = fmt.Sprintf("Kolom %s harus berupa alamat email yang valid", field)
			case "len":
				errs[field] = fmt.Sprintf("Kolom %s harus berukuran tepat %s karakter", field, e.Param())
			case "oneof":
				errs[field] = fmt.Sprintf("Kolom %s harus salah satu dari: %s", field, e.Param())
			case "uuid":
				errs[field] = fmt.Sprintf("Kolom %s harus berformat UUID yang valid", field)
			default:
				errs[field] = fmt.Sprintf("Kolom %s tidak valid", field)
			}
		}
		return errs
	}
	
	// If it's not a ValidationErrors (e.g. malformed JSON)
	errs["request"] = "Format data JSON tidak valid"
	return errs
}
