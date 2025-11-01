package helpers

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func TranslateErrorMessage(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validatioErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatioErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				errorsMap[field] = fmt.Sprintf("%s is required", field)
			case "email":
				errorsMap[field] = "invalid email format"
			case "unique":
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			case "min":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "max":
				errorsMap[field] = fmt.Sprintf("%s must be at least %s characters", field, fieldError.Param())
			case "numeric":
				errorsMap[field] = fmt.Sprintf("%s must be a number", field)
			default:
				errorsMap[field] = "invalid value"
			}
		}
	}

	if err != nil {
		if strings.Contains(err.Error(), "not the hash of the given password") {
			errorsMap["Error"] = "Password doesn't match"
		}

		if strings.Contains(err.Error(), "Duplicate entry") {
			field := extractDuplicateField(err.Error())
			if field == "" {
				errorsMap[field] = fmt.Sprintf("%s already exists", field)
			} else {
				errorsMap["Error"] = "Duplicate entry"
			}
		} else if err == gorm.ErrRecordNotFound {
			errorsMap["Error"] = "Record not found"
		}
	}

	return errorsMap
}

func extractDuplicateField(errMsg string) string {
	// Contoh error MySQL: Error 1062: Duplicate entry 'test@example.com' for key 'users.email'
	// Kita ambil bagian setelah 'for key' lalu extract nama field
	re := regexp.MustCompile(`for key '(\w+\.)?(\w+)'`)
	matches := re.FindStringSubmatch(errMsg)
	if len(matches) == 3 {
		// Hasilkan kapitalisasi nama field
		return strings.Title(matches[2])
	}
	return ""
}
