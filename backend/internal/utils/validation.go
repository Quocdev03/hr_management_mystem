package utils

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"unicode/utf8"
)

// --- Constants for Field Names ---
const (
	FieldEmail        = "email"
	FieldPassword     = "password"
	FieldUserName     = "user_name"
	FieldPhone        = "phone"
	FieldFirstName    = "first_name"
	FieldLastName     = "last_name"
	FieldCode         = "code"
	FieldName         = "name"
	FieldDepartmentID = "department_id"
	FieldPosition     = "position"
	FieldSalary       = "salary"
	FieldJoinDate     = "join_date"
	FieldBirthDate    = "birth_date"
	FieldGender       = "gender"
	FieldStatus       = "status"
	FieldDescription  = "description"
	FieldManagerID    = "manager_id"
	FieldRoleID       = "role_id"
	FieldPage         = "page"
	FieldLimit        = "limit"
)

// --- Compiled Regexes ---
var (
	regUpper = regexp.MustCompile(`[A-Z]`)
	regLower = regexp.MustCompile(`[a-z]`)
	regDigit = regexp.MustCompile(`[0-9]`)
)

// FieldError chứa thông tin lỗi cho từng field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors là danh sách các lỗi validation
type ValidationErrors struct {
	Errors []FieldError `json:"errors"`
}

func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, FieldError{
		Field:   field,
		Message: message,
	})
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return ""
	}
	var msgs []string
	for _, e := range ve.Errors {
		msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field, e.Message))
	}
	return strings.Join(msgs, "; ")
}

// --- Helper methods ---

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func length(s string) int {
	return utf8.RuneCountInString(s)
}

// --- Reusable Validation Checkers ---

func validateEmailFormat(ve *ValidationErrors, email string) {
	if _, err := mail.ParseAddress(email); err != nil {
		ve.Add(FieldEmail, "Email không đúng định dạng")
	}
}

func CheckEmail(ve *ValidationErrors, email string) {
	email = normalize(email)
	if email == "" {
		ve.Add(FieldEmail, "Email là bắt buộc")
		return
	}
	validateEmailFormat(ve, email)
}

func validatePasswordFormat(ve *ValidationErrors, password string) {
	if len(password) < 8 {
		ve.Add(FieldPassword, "Mật khẩu phải có ít nhất 8 ký tự")
		return
	}
	if !regUpper.MatchString(password) ||
		!regLower.MatchString(password) ||
		!regDigit.MatchString(password) {
		ve.Add(FieldPassword, "Mật khẩu phải chứa ít nhất 1 chữ hoa, 1 chữ thường và 1 số")
	}
}

func CheckPassword(ve *ValidationErrors, password string) {
	if normalize(password) == "" {
		ve.Add(FieldPassword, "Mật khẩu là bắt buộc")
		return
	}
	validatePasswordFormat(ve, password)
}
