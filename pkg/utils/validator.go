package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// สร้าง Instance เดียวใช้ทั่วโปรเจค
var validate = validator.New()

// Struct สำหรับเก็บ Error ที่อ่านง่าย
type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

// ฟังก์ชัน Validate Struct
func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.Field() // ชื่อ Field (เช่น Email)
			element.Tag = err.Tag()     // ชื่อ Tag ที่ผิด (เช่น required, email)
			element.Value = err.Param() // ค่า parameter (เช่น min=8)
			errors = append(errors, &element)
		}
		return errors
	}

	return nil
}

// ฟังก์ชันช่วยจัด Format Error เป็น String เดียว (เผื่อไม่อยากส่งเป็น Array)
func FormatValidationError(errors []*ErrorResponse) error {
	if len(errors) == 0 {
		return nil
	}
	var errMsgs []string
	for _, err := range errors {
		errMsgs = append(errMsgs, fmt.Sprintf("[%s]: failed on tag '%s'", err.Field, err.Tag))
	}
	return fmt.Errorf("%s", strings.Join(errMsgs, ", "))
}
